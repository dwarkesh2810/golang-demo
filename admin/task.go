package admin

// The first 6 parts are:
//       second: 0-59
//       minute: 0-59
//       hour: 1-23
//       day: 1-31
//       month: 1-12
//       weekdays: 0-6（0 is Sunday）

// Some special sign:
//
//	*: any time
//	,: separator. E.g.: 2,4 in the third part means run at 2 and 4 o'clock
//
// 　　    －: range. E.g.: 1-5 in the third part means run between 1 and 5 o'clock
//
//	/n : run once every n time. E.g.: */1 in the third part means run once every an hour. Same as 1-23/1
//
// ///////////////////////////////////////////////////////
//
//	0/30 * * * * *                        run every 30 seconds
//	0 43 21 * * *                         run at 21:43
//	0 15 05 * * *                         run at 05:15
//	0 0 17 * * *                          run at 17:00
//	0 0 17 * * 1                          run at 17:00 of every Monday
//	0 0,10 17 * * 0,2,3                   run at 17:00 and 17:10 of every Sunday, Tuesday and Wednesday
//	0 0-10 17 1 * *                       run once every minute from 17:00 to 7:10 on 1st day of every month
//	0 0 0 1,15 * 1                        run at 0:00 on 1st and 15th of each month and every Monday
//	0 42 4 1 * *                          run at 4:42 on 1st of every month
//	0 0 21 * * 1-6                        run at 21:00 from Monday to Saturday
//	0 0,10,20,30,40,50 * * * *            run every 10 minutes
//	0 */10 * * * *                        run every 10 minutes
//	0 * 1 * * *                           run every one minute from 1:00 to 1:59
//	0 0 1 * * *                           run at 1:00
//	0 0 */1 * * *                         run at :00 of every hour
//	0 0 * * * *                           run at :00 of every hour
//	0 2 8-20/3 * * *                      run at 8:02, 11:02, 14:02, 17:02 and 20:02
//	0 30 5 1,15 * *                       run at 5:30 of 1st and 15th of every month
import (
	"context"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/task"
	"github.com/dwarkesh2810/golang-demo/helpers"
	"github.com/dwarkesh2810/golang-demo/models"
)

func CreateTask(taskName, schedule string, f task.TaskFunc) {
	tasks := task.NewTask(taskName, schedule, f)
	task.AddTask(taskName, tasks)
	task.StartTask()
}


func SendPendingEmail(c context.Context) error {
	o := orm.NewOrm()
	// orm.Debug = true
	var emails []models.EmailLogs
	_, err := o.QueryTable(new(models.EmailLogs)).Filter("status", "pending").All(&emails, "LogId", "emailTo", "name", "subject", "body")
	if err != nil {
		return err
	}
	for _, email := range emails {
		sent, _ := helpers.SendMailOTp(email.To, email.Name, email.Subject, email.Body)
		if sent {
			var UpdateEmail = models.EmailLogs{Id: email.Id, Status: "success"}
			_, err := o.Update(&UpdateEmail, "status")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
