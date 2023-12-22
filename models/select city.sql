select city.city_name as city_name,city.city_id as city_id from countries as c
         LEFT JOIN cities as city ON city.country_id=c.country_id
         LEFT JOIN states as s ON s.state_id = city.state_id
         WHERE upper(c.country_name) = upper('India') 
AND upper(s.state_name)=upper('GUjarat') 
or upper(c.country_name) = upper('') 
or  upper(s.state_name)=upper('')
