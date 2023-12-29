package helpers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/dwarkesh2810/golang-demo/conf"
	"github.com/dwarkesh2810/golang-demo/dto"

	beego "github.com/beego/beego/v2/server/web"

	"github.com/beego/beego/v2/server/web/context"

	"github.com/dgrijalva/jwt-go"
	"github.com/jung-kurt/gofpdf"
	"github.com/tealeg/xlsx"
	"golang.org/x/crypto/bcrypt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/i18n"
	"github.com/go-ini/ini"
)

func SecondsToDayHourMinAndSeconds(seconds int64) (int64, int64, int64, int64) {
	days := seconds / 86400
	hour := (seconds % 86400) / 3600
	minute := (seconds % 3600) / 60
	second := seconds % 60
	return int64(days), int64(hour), int64(minute), int64(second)
}

func HashData(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
func VerifyHashedData(hashedString string, dataString string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(dataString))
	if err != nil {
		return errors.New("HASHED_ERROR")
	}
	return nil
}

// >>>>>>>>>>>>>Function to send mail <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<,
func SendMailOTp(userEmail string, name string, subject string, body string) (bool, error) {
	from := conf.Env.FromEmail
	password := conf.Env.Password
	to := []string{
		userEmail,
	}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte("Subject: " + subject + "\r\n" + mime + "\r\n" + body)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return false, errors.New("TWILIO_ERROR")
	}
	return true, nil
}

/*PAGINATION FUNCTION PROVIDE ALL DETAILS LIKE AS  CURRENT PAGE,LAST PAGE AND TOTAL ROWS AND TOTAL PAGES IT ALSO */
func Pagination(current_page, pageSize int, tableName string, totalMatchCount int) (map[string]interface{}, error) {
	db := orm.NewOrm()
	if current_page <= 0 {
		current_page = 1
	}
	if pageSize < 0 {
		pageSize = 10
	}

	var totalRows int
	if totalMatchCount > 0 {
		totalRows = totalMatchCount
	}

	if totalMatchCount == 0 && tableName != "" {
		err := db.Raw(`SELECT COUNT(*) as totalRows FROM ` + tableName).QueryRow(&totalRows)
		if err != nil {
			return nil, err
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(pageSize)))
	lastPageNumber := totalPages
	if lastPageNumber == 0 {
		lastPageNumber = 1
	}

	previousPageNumber := current_page - 1
	if previousPageNumber < 1 {
		previousPageNumber = 0
	}

	nextPageNumber := current_page + 1
	if nextPageNumber > totalPages {
		nextPageNumber = totalPages
	}

	pagination_data := map[string]interface{}{
		"CurrentPage":   current_page,
		"PreviousPage":  previousPageNumber,
		"NextPage":      nextPageNumber,
		"PerPageRecord": pageSize,
		"TotalRows":     totalRows,
		"TotalPages":    totalPages,
		"LastPage":      lastPageNumber,
	}
	if current_page > lastPageNumber {
		pagination_data["pageOpen_error"] = 1
		pagination_data["current_page"] = current_page
		pagination_data["last_page"] = lastPageNumber
	}

	return pagination_data, nil
}

/*
FORMAT DATE TIME FUNCTION TAKE DATE LIKE [2023-12-11 10:11:38.804636+05:30]
AND IF RETURNTYPE NOT PASS THAN IT  RETURNS DATE AND TIME DATE:- DD-MM-YY AND ALSO RETURNS
IF PASS TIME THAN IT RETURNS FORMAT:-  HH:MM:SS AM/PM
*/

func FormatDateTime(inputDateTime string, formatType ...string) (map[string]string, error) {
	inputLayouts := []string{
		"2006-01-02 15:04:05.999999999 -0700 MST",
		"2006-01-02 15:04:05.999999-07",
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02 15:04:05.999999999 -0700 MST",
		"2006-01-02 15:04:05.999999999 -0700 MST",
		"2006-01-02 15:04:05.999999-07",
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02 15:04:05.999999999 -0700 MST m=+0",
		"2006-01-02 15:04:05.999999999 -0700 MST m=+0.000000001",
		"2006-01-02 15:04:05.999999999 -0700 MST m=+0.000000001",
	}

	var parsedTime time.Time
	var err error

	for _, layout := range inputLayouts {
		parsedTime, err = time.Parse(layout, inputDateTime)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	dateLayoutDefault := "02-01-2006"
	dateLayoutISO := "2006-01-02"
	timeLayout := "03:04:05 PM"

	dayLayout := parsedTime.Format("Monday")

	result := make(map[string]string)

	switch len(formatType) {
	case 0:
		result["date"] = parsedTime.Format(dateLayoutDefault)
		result["time"] = parsedTime.Format(timeLayout)
		result["day"] = dayLayout
	case 1:
		switch strings.ToUpper(formatType[0]) {
		case "DATE":
			result["date"] = parsedTime.Format(dateLayoutISO)
			result["time"] = parsedTime.Format(timeLayout)
			result["day"] = dayLayout
		case "TIME":
			result["date"] = parsedTime.Format(dateLayoutDefault)
			result["time"] = parsedTime.Format(timeLayout)
			result["day"] = dayLayout
		case "DAY":
			result["date"] = parsedTime.Format(dateLayoutDefault)
			result["time"] = parsedTime.Format(timeLayout)
			result["day"] = parsedTime.Format(dayLayout)
		case "DIFF":
			currentTime := time.Now()
			difference := currentTime.Sub(parsedTime).Hours() / 24
			result["date"] = parsedTime.Format(dateLayoutDefault)
			result["time"] = parsedTime.Format(timeLayout)
			result["day"] = parsedTime.Format(dayLayout)
			result["diff"] = fmt.Sprintf(" %.f days", difference)
		default:
			return nil, fmt.Errorf("unsupported format type")
		}
	default:
		return nil, fmt.Errorf("too many arguments")
	}

	return result, nil
}

func GetFormatedDate(date time.Time, formate string) string {
	var formatedDate string
	switch formate {
	case "dd-mm-yy":
		inputTime := date
		day, month, year := inputTime.Day(), inputTime.Month(), inputTime.Year()%100
		formatedDate = fmt.Sprintf("%02d-%02d-%d", day, month, year)
	case "dd-mm-yyyy":
		inputTime := date
		day, month, year := inputTime.Day(), inputTime.Month(), inputTime.Year()
		formatedDate = fmt.Sprintf("%02d-%02d-%d", day, month, year)
	case "yyyy-mm-dd":
		inputTime := date
		day, month, year := inputTime.Day(), inputTime.Month(), inputTime.Year()
		formatedDate = fmt.Sprintf("%02d-%02d-%d", year, month, day)
	case "mm-dd-yyyy":
		inputTime := date
		day, month, year := inputTime.Day(), inputTime.Month(), inputTime.Year()
		formatedDate = fmt.Sprintf("%02d-%02d-%d", year, month, day)
	case "dd-mm":
		inputTime := date
		day, month := inputTime.Day(), inputTime.Month()
		formatedDate = fmt.Sprintf("%02d-%02d", day, month)
	default:
		formatedDate = "not formated"
	}
	return formatedDate
}

/*END FORMATE DATE TIME FUNCTION*/

/*----------LANGUAGE TRANSLATION FUNCTION START-----------------*/
var defaultLang = "en-US"

func Init() {
	beego.InsertFilter("*", beego.BeforeRouter, func(ctx *context.Context) {
		lang := GetLanguageFromMultipleSources(ctx)
		SetLanguage(ctx, lang)
	})
	beego.InsertFilter("*", beego.AfterExec, func(ctx *context.Context) {

	})
}
func GetLanguageFromMultipleSources(ctx *context.Context) string {
	if lang := ctx.Input.Query("lang"); lang != "" && isValidLanguage(lang) {
		return lang
	}
	if lang := ctx.Input.Header("lang"); lang != "" && isValidLanguage(lang) {
		return lang
	}
	if lang := ctx.Input.Cookie("lang"); lang != "" && isValidLanguage(lang) {
		return lang
	}
	return "en-US"
}
func isValidLanguage(lang string) bool {
	lang = strings.ToUpper(lang)
	allowedLanguages := map[string]bool{"EN-US": true, "EN-GB": true, "HI-IN": true, "BO-IN": true, "EN-IN": true, "PS-AF": true, "GU-IN": true, "KN-IN": true, "MR-IN": true, "NE-IN": true, "OR-IN": true, "TA-IN": true, "TE-IN": true, "UR-IN": true, "FR-FR": true, "RU-RU": true, "IG-NG": true}
	return allowedLanguages[lang]
}
func SetLanguage(ctx *context.Context, lang string) {
	ctx.Input.SetData("lang", lang)
	i18n.SetMessageWithDesc(lang, "conf/language/locale_"+lang+".ini", "conf/language/locale_"+lang+".ini")

	ctx.SetCookie("lang", lang, 24*60*60, "/")

	defaultLang = lang
}
func Translate(ctx *context.Context, key string) string {
	langKey := GetLanguageFromMultipleSources(ctx)
	langTrans := strings.Split(langKey, "-")
	langTrans[0] = strings.ToLower(langTrans[0])
	if len(langTrans) > 1 {
		langTrans[1] = strings.ToUpper(langTrans[1])
	}
	langKey = strings.Join(langTrans, "-")
	SetLanguage(ctx, langKey)
	return i18n.Tr(defaultLang, key)
}
func TranslateMessage(ctx *context.Context, section, sectionMessage string) string {

	translationKey := fmt.Sprintf("%s.%s", section, sectionMessage)
	return Translate(ctx, translationKey)
}

/*CREATE INI FILE ACCORDING TO LANGUAGE CODE  IN DIRECTORY OF [CONF/LANGUAGE]/*/
func CreateINIFiles(data []map[string]string) error {
	for _, item := range data {
		languageCode := item["language_code"]

		fileName := fmt.Sprintf("locale_%s.ini", languageCode)
		filePath := filepath.Join("conf/language", fileName)

		if err := os.MkdirAll("conf/language", os.ModePerm); err != nil {
			return err
		}
		cfg, err := ini.Load(filePath)
		if err != nil {
			cfg = ini.Empty()
		}

		section, err := cfg.NewSection(item["section"])
		if err != nil {
			return err
		}

		for key := range item {
			_ = key
			section.NewKey(item["lable_code"], item["language_value"])
		}

		err = cfg.SaveTo(filePath)
		if err != nil {
			return err
		}

	}

	return nil
}
func formateCSVDate(value interface{}) string {
	switch v := value.(type) {
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprintf("%v", value)
	}
}
func formatValue(value interface{}) interface{} {
	switch v := value.(type) {
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	default:
		return v
	}
}
func sumSliceElements(slice []float64) float64 {
	var total float64
	for _, value := range slice {
		total += value
	}
	return total
}

/*CREATE FILE [XLSX,PDF,CSV] IN SPECIFIC DIRECTORY*/
func CreateFile(data []map[string]interface{}, headers []string, folderPath, fileNamePrefix, fileType string) (string, error) {
	TYPE := strings.ToUpper(fileType)
	switch TYPE {
	case "XLSX":
		file := excelize.NewFile()
		sheet := "Sheet1"
		file.NewSheet(sheet)

		// Set header row
		for colNum, header := range headers {
			cell := fmt.Sprintf("%c%d", 'A'+colNum, 1)
			file.SetCellValue(sheet, cell, header)
		}

		// Set data rows
		for rowNum, rowData := range data {
			for colNum, key := range headers {
				cell := fmt.Sprintf("%c%d", 'A'+colNum, rowNum+2)
				if value, ok := rowData[key]; ok {
					file.SetCellValue(sheet, cell, formatValue(value))
				}
			}
		}

		if folderPath == "" {
			folderPath = "assets/uploads/FILES/XLSX"
		}

		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create folder: %v", err)
		}

		fileName := fmt.Sprintf("%s_%s.xlsx", fileNamePrefix, time.Now().Format("20060102150405"))
		filePath := filepath.Join(folderPath, fileName)
		if err := file.SaveAs(filePath); err != nil {
			return "", err
		}
		return filePath, nil

	case "CSV":
		if folderPath == "" {
			folderPath = "assets/uploads/FILES/CSV"
		}

		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create folder: %v", err)
		}

		fileName := fmt.Sprintf("%s_%s.csv", fileNamePrefix, time.Now().Format("20060102150405"))
		filePath := filepath.Join(folderPath, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to create CSV file: %v", err)
		}
		defer file.Close()

		csvWriter := csv.NewWriter(file)
		defer csvWriter.Flush()

		// Write header row
		if err := csvWriter.Write(headers); err != nil {
			return "", fmt.Errorf("failed to write CSV header: %v", err)
		}

		// Write data rows
		for _, rowData := range data {
			var row []string
			for _, key := range headers {
				if value, ok := rowData[key]; ok {
					row = append(row, formateCSVDate(value))
				} else {
					row = append(row, "") // Handle missing data
				}
			}
			if err := csvWriter.Write(row); err != nil {
				return "", fmt.Errorf("failed to write CSV row: %v", err)
			}
		}

		return filePath, nil

	case "PDF":
		pdf := gofpdf.New("L", "mm", "A4", "")
		pdf.AddPage()
		fontSize := 10.0
		pdf.SetFont("Arial", "B", fontSize)

		pageWidth, _ := pdf.GetPageSize()

		colWidths := make([]float64, len(headers))
		totalWidth := pageWidth - 20
		for colNum, header := range headers {
			colWidths[colNum] = pdf.GetStringWidth(header) + 6
		}

		scaleFactor := totalWidth / sumSliceElements(colWidths)
		for colNum := range colWidths {
			colWidths[colNum] *= scaleFactor
		}

		for colNum, header := range headers {
			pdf.CellFormat(colWidths[colNum], 10, header, "1", 0, "", false, 0, "")
		}

		pdf.Ln(-1)

		pdf.SetFont("Arial", "", fontSize)

		for _, rowData := range data {

			for colNum, key := range headers {
				if value, ok := rowData[key]; ok {
					cellValue := fmt.Sprintf("%v", formatValue(value))
					pdf.CellFormat(colWidths[colNum], 10, cellValue, "1", 0, "", false, 0, "")
				}
			}

			pdf.Ln(-1)
		}

		// If filepath not given, it takes the default
		if folderPath == "" {
			folderPath = "assets/uploads/FILES/PDF"
		}

		// If the folder is not present in the directory, it creates a new folder directory
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create folder: %v", err)
		}

		// Generate file name
		fileName := fmt.Sprintf("%s_%s.pdf", fileNamePrefix, time.Now().Format("20060102150405"))
		filePath := filepath.Join(folderPath, fileName)

		// Save the PDF file
		if err := pdf.OutputFileAndClose(filePath); err != nil {
			return "", err
		}

		return filePath, nil

	default:
		return "", fmt.Errorf("unsupported file type: %s", fileType)
	}
}

/* ----------------------end XLSX file creating functions---------------------------------------------------------*/

/*-------------------------------XLSX AND CSV FILE READING FUNCTION*/

func ReadXLSXFiles(filePath string) ([]map[string]interface{}, error) {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	var allRows []map[string]interface{}

	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			rowData := make(map[string]interface{})
			for index, cell := range row.Cells {
				rowData[fmt.Sprintf("Column%d", index+1)] = cell.String()
			}

			allRows = append(allRows, rowData)
		}
	}
	return allRows, nil
}

func ReadXLSXFile(filePath string) ([]map[string]interface{}, error) {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	var allRows []map[string]interface{}

	// Assume the first row is the header
	var headerRow []string
	if len(xlFile.Sheets) > 0 && len(xlFile.Sheets[0].Rows) > 0 {
		headerRow = make([]string, len(xlFile.Sheets[0].Rows[0].Cells))
		for index, cell := range xlFile.Sheets[0].Rows[0].Cells {
			headerRow[index] = cell.String()
		}
	}

	for _, sheet := range xlFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			if rowIndex == 0 {
				// Skip the header row
				continue
			}

			rowData := make(map[string]interface{})
			for index, cell := range row.Cells {
				if index < len(headerRow) {
					rowData[headerRow[index]] = cell.String()
				}
			}

			allRows = append(allRows, rowData)
		}
	}

	return allRows, nil
}

func ReadCSVFile(filePath string) ([]map[string]interface{}, error) {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	var allRows []map[string]interface{}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	columnHeaders := records[0]

	for _, dataRow := range records[1:] {
		rowData := make(map[string]interface{})

		for index, value := range dataRow {
			rowData[columnHeaders[index]] = value
		}

		allRows = append(allRows, rowData)
	}

	return allRows, nil
}

/* END XLSX AND CSV FILE READING FUNCTION END-----------------------*/

/*GENERATE UNIQUE CODE (ALPH + NUMERIC) USE COMBINATION ACCORDING YOUR REQUIREMENT OF CODE LENGTH HERE ONLY YOU PASS LENGTH FUNCTION GENERATE UNIQUE CODE*/
func GenerateUniqueCodeString(length int) string {
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

/*HASH PASSWORD*/
func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

/*VERIFY HASH PASSWORD*/
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

/*GET TOKEN SET CLAIMS */
func GetTokenClaims(c *context.Context) map[string]interface{} {
	token_claims := c.Input.GetData("user")
	user_id := token_claims.(jwt.MapClaims)["ID"]
	user_email := token_claims.(jwt.MapClaims)["Email"]
	response := map[string]interface{}{"User_id": user_id, "User_Email": user_email}
	return response
}

/*UPLOAD FILE ACCORDING TO THE UPLOAD DIRECTORY PATH*/
func UploadFile(fileToUpload multipart.File, fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	defer fileToUpload.Close()
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	if err := os.MkdirAll(uploadDir, 0777); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}
	filePath := filepath.Join(uploadDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create the destination file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, fileToUpload)
	if err != nil {
		return "", fmt.Errorf("failed to copy the file: %v", err)
	}

	return filePath, nil
}

/*REMOVE FILE BY USING FILE NAME AND DIRECTORY*/
func RemoveFile(fileName, directory string) error {
	err := os.Remove(filepath.Join(directory, fileName))
	if err != nil {
		return err
	}
	return nil
}

/*REMOVE FILE BY THE FILE PATH*/
func RemoveFileByPath(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

/*SPLITE FILE PATH FROM THE LAST /(SLASH) */
func SplitFilePath(SplitString string) (string, string) {
	lastIndex := strings.LastIndex(SplitString, "/")

	var fileDirectory string
	var fileName string

	if lastIndex != -1 {
		fileDirectory = SplitString[:lastIndex]
		fileName = SplitString[lastIndex+1:]
	} else {
		fileDirectory = "No '/' found in the string."
		fileName = fileDirectory
	}

	return fileName, fileDirectory
}

/*GENERATE UNIQUE CODE WITH UNDERSCORE AFTER WITHSTRING EX. dev_12*/
func UniqueCode(number int, withString string) string {
	withString = strings.ReplaceAll(withString, " ", "_")
	result := fmt.Sprintf("%s_%d", withString, number)
	return strings.ToUpper(result)
}

/*EXTRACT KEYS FROM THE []MAP[STRING]INTERFACE{} AND CONVERT INTO []STRING*/
func ExtractKeys(data []map[string]interface{}) []string {
	keys := make(map[string]struct{})
	for _, item := range data {
		for key := range item {
			keys[key] = struct{}{}
		}
	}
	var result []string
	for key := range keys {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

/************************************** CONVERTIONS FUNCTIONS *********************************************************/

/*CONVERT ORM.PARMS TO []MAP[STRING]STRING MAP SLICE FORMAT*/
func ConvertToMapSlice(results []orm.Params) ([]map[string]string, error) {
	var converted []map[string]string
	for _, params := range results {
		convertedItem := make(map[string]string)
		for key, value := range params {
			convertedItem[key] = fmt.Sprintf("%v", value)
		}
		converted = append(converted, convertedItem)
	}
	return converted, nil
}

/*TRANSLATE DATA INTO KEY VALUE PAIRS IT WILL WORK BOTH []ORM.PARAMS AND  DATA INTERFACE{}*/
func TransformToKeyValuePairs(data interface{}) ([]map[string]interface{}, error) {
	value := reflect.ValueOf(data)

	if value.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input data must be a slice")
	}

	result := make([]map[string]interface{}, value.Len())

	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		fields := make(map[string]interface{})

		switch item.Interface().(type) {
		case orm.Params:
			for key, value := range item.Interface().(orm.Params) {
				fields[key] = value
			}
		case map[string]interface{}:
			fields = item.Interface().(map[string]interface{})
		default:
			if item.Kind() == reflect.Struct {
				for j := 0; j < item.NumField(); j++ {
					field := item.Type().Field(j)
					fieldName := field.Tag.Get("json")
					if fieldName == "" {
						fieldName = strings.ToLower(field.Name)
					}
					fields[fieldName] = item.Field(j).Interface()
				}
			} else {
				return nil, fmt.Errorf("items in the slice must be orm.Params, map[string]interface{}, or structs")
			}
		}

		result[i] = fields
	}

	return result, nil
}

/*FOR TRANSLATE IN LANGUAGE WHICH IS ALREADY SET IN CONTEXT/COOKIE/HEADER*/

func LanguageTranslate(c beego.Controller, key string) string {
	lang := c.Ctx.Input.GetData("lang").(string)
	language := strings.ToLower(lang)
	switch language {
	case "en-us":
		lang = "en-US"
	case "hi-in":
		lang = "hi-IN"
	case "zh-cn":
		lang = "zh-CN"
	}
	return i18n.Tr(lang, key)
}

/*FOR CHECKING ELEMENT IN SLICE*/
func CheckIfExists(elemet string, data []string) bool {
	for i := 0; i < len(data); i++ {
		if strings.EqualFold(elemet, data[i]) {
			return true
		}
	}
	return false
}

/*Function for checck car typy is valid or not*/
func NewCarType(input string) (dto.CarType, error) {
	switch input {
	case "sedan", "hatchback", "SUV":
		return dto.CarType(input), nil
	default:
		return "", errors.New("INVALID_CAR")
	}
}

/*FILTER DATA ACCORDING TO QUERY AND GIVE FILTER DATA COUNTS*/
func FilterData(currentPage, pageSize int, query, tableName string, searchFields map[string]string, applyPosition, countQuery string, otherFieldSCount int) ([]orm.Params, map[string]interface{}, int, error) {
	db := orm.NewOrm()
	if currentPage <= 0 {
		currentPage = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (currentPage - 1) * pageSize
	var homeResponse []orm.Params

	if len(searchFields) > 0 {
		whereClause := generateWhereClause(searchFields, applyPosition, otherFieldSCount)
		query += " " + whereClause
		countQuery += " " + whereClause
	}
	query += `
        LIMIT ? OFFSET ?
    `
	_, err := db.Raw(query, append(generateSearchParameters(searchFields, applyPosition), pageSize, offset)...).Values(&homeResponse)
	if err != nil {
		return nil, nil, 0, err
	}
	var totalMatchData int
	err = db.Raw(countQuery, generateSearchParameters(searchFields, applyPosition)...).QueryRow(&totalMatchData)
	if err != nil {
		return nil, nil, 0, err
	}

	paginationData, paginationErr := Pagination(currentPage, pageSize, tableName, totalMatchData)
	if paginationErr != nil {
		return nil, paginationData, 0, paginationErr
	}
	return homeResponse, paginationData, totalMatchData, nil
}

func generateWhereClause(fields map[string]string, applyPosition string, otherFieldCount int) string {
	var conditions []string
	applyPosition = strings.ToUpper(applyPosition)
	for field, value := range fields {
		if value != "" {
			condition := ""
			if applyPosition == "" {
				condition = field + " ILIKE ?"
			} else if applyPosition == "START" {
				condition = field + " ILIKE ?"
			} else {
				condition = field + " ILIKE ?"
			}
			conditions = append(conditions, condition)
		}
	}
	var otherFileds int = 0
	if otherFieldCount > 0 {
		otherFileds = otherFieldCount
	}
	if otherFileds > 0 && len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " OR ")
		return whereClause
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " OR ")
		return whereClause
	}
	return ""
}

func generateSearchParameters(fields map[string]string, applyPostion string) []interface{} {
	var parameters []interface{}
	applyPostion = strings.ToUpper(applyPostion)
	for _, field := range fields {
		if field != "" {
			if applyPostion == "" {
				parameters = append(parameters, "%"+field+"%")
			} else if applyPostion == "START" {
				parameters = append(parameters, field+"%")
			} else {
				parameters = append(parameters, "%"+field)
			}
		}
	}

	return parameters
}

/*END FILTER DATA FUNCTION*/

/*RETURN RECORDS WITH PAGINATION*/
func FetchDataWithPaginations(current_page, pageSize int, tableName, query string) ([]orm.Params, map[string]interface{}, error) {
	db := orm.NewOrm()
	if current_page <= 0 {
		current_page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (current_page - 1) * pageSize

	var homeResponse []orm.Params
	_, err := db.Raw(query, pageSize, offset).Values(&homeResponse)

	if err != nil {
		return nil, nil, err
	}
	pagination_data, pagination_err := Pagination(current_page, pageSize, tableName, 0)
	if pagination_err != nil {
		return nil, pagination_data, nil
	}
	pagination_data["matchCount"] = 0
	return homeResponse, pagination_data, nil
}

/*CONVERT STRUCT TO MAP*/
func ConvertStructToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct or a pointer to struct")
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		result[fieldName] = field.Interface()
	}

	return result, nil
}

/*-----------------------Image validation----------------------------------*/
func GetFileExtension(file string) string {
	splitFileName := strings.Split(file, ".")
	return strings.ToLower(splitFileName[len(splitFileName)-1])
}

/*-------------------------------end----------------------------------------*/

/*------------------------------------Pagination for search ------------------------------------------------*/

func PaginationForSearch(current_page, pageSize int, totalRecordQuery, matchCountQuery, mainRecordQuery string) ([]orm.Params, map[string]interface{}, error) {
	db := orm.NewOrm()
	// orm.Debug = true
	if pageSize <= 0 {
		pageSize = 10
	}
	if current_page <= 0 {
		current_page = 1
	}
	offset := (current_page - 1) * pageSize
	var states []orm.Params
	_, err := db.Raw(mainRecordQuery, pageSize, offset).Values(&states)
	if err != nil {
		return nil, nil, err
	}
	var totalRows int
	err = db.Raw(totalRecordQuery).QueryRow(&totalRows)
	if err != nil {
		return nil, nil, err
	}
	var count []orm.Params
	_, err = db.Raw(matchCountQuery).Values(&count)
	if err != nil {
		return nil, nil, err
	}
	totalPages := int(math.Ceil(float64(len(count)) / float64(pageSize)))

	lastPageNumber := totalPages
	if lastPageNumber == 0 {
		lastPageNumber = 1
	}
	previousPageNumber := current_page - 1
	if previousPageNumber < 1 {
		previousPageNumber = 0
	}
	nextPageNumber := current_page + 1
	if nextPageNumber > totalPages {
		nextPageNumber = totalPages
	}
	pagination_data := map[string]interface{}{
		"CurrentPage":   current_page,
		"PreviousPage":  previousPageNumber,
		"NextPage":      nextPageNumber,
		"PerPageRecord": pageSize,
		"TotalRows":     totalRows,
		"TotalPages":    totalPages,
		"LastPage":      lastPageNumber,
	}
	if current_page > lastPageNumber {
		pagination_data["pageOpen_error"] = 1
		pagination_data["current_page"] = current_page
		pagination_data["last_page"] = lastPageNumber
	}
	pagination_data["matchCount"] = len(count)
	return states, pagination_data, nil
}

func CapitalizeWords(s string) string {
	words := strings.Fields(s) // Split the string into words
	capitalizedWords := make([]string, len(words))

	for i, word := range words {
		if len(word) > 0 { // Ensure the word is not empty
			capitalizedWords[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return strings.Join(capitalizedWords, " ") // Join the capitalized words back into a string
}

/*VALID FORMAT CONVERTIONS FOR INI LANGUAGE CODE*/
func ConvertIntoIniFormateCode(input string) string {
	regex := regexp.MustCompile("[^a-zA-Z0-9]")
	input = regex.ReplaceAllString(input, "-")
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return input
	}
	return strings.ToLower(parts[0]) + "-" + strings.ToUpper(parts[1])
}

/*READ INI FILE AND CONVERT SECTION AND KEY VALUE PAIRS IN MAP[STRING]MAP[STRING]STRING*/
func ParseINIFile(filePath string) (map[string]map[string]string, error) {
	config, err := ini.Load(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load INI file: %w", err)
	}
	dataMap := make(map[string]map[string]string)

	for _, section := range config.Sections() {
		if section.Name() == "DEFAULT" {
			continue // Skip the DEFAULT section
		}

		dataMap[section.Name()] = make(map[string]string)

		for _, key := range section.Keys() {
			dataMap[section.Name()][key.Name()] = key.String()
		}
	}

	return dataMap, nil
}

/*CONVERT INTO MAP FOR PROCESSMAP DATA FUNCTIONS*/
func ConvertToDataMap(inputMap map[string]map[string]string) map[string]map[string]map[string]string {
	dataMap := make(map[string]map[string]map[string]string)
	for section, sectionMap := range inputMap {
		sectionDataMap := make(map[string]map[string]string)

		for k, v := range sectionMap {
			sectionDataMap[k] = map[string]string{
				"language_value": v,
			}
		}
		dataMap[section] = sectionDataMap
	}

	return dataMap
}

/*EXTRACT INI LANGUAGE CODE FROM FILE EX.locale_mu-IN.ini [OUTPUT STRING mu-IN]*/
func ExtractLanguageCode(fileName string) string {
	languageCode := strings.TrimPrefix(fileName, "locale_")
	languageCode = strings.TrimSuffix(languageCode, ".ini")

	return languageCode
}
