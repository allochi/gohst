-- Just select the json field
select "data" from "DataStore";

-- Select all fields in rows and convert them to json with proper field names
select row_to_json(contacts) from contacts;

-- Select all few fields in all rows and convert them to json but field names are lost
select row_to_json(row(id, first_name, last_name)) from contacts;

-- Select all rows in a query with specific field and turn them to json with proper field names
select row_to_json(row_data) from (
  select c.id, c.first_name, c.last_name, j.name as job_title  from contacts as c left join job_titles as j on c.job_title_id = j.id
) row_data;

-- Convert the json rows in one big array in one row
select array_to_json(array_agg(row_to_json(row_data))) from (
  select c.id, c.first_name, c.last_name, j.name as job_title  from contacts as c left join job_titles as j on c.job_title_id = j.id
) row_data;

CREATE or REPLACE FUNCTION getTax(subtotal REAL) RETURNS REAL AS $$
BEGIN
  RETURN subtotal * 0.06;
END;
$$ LANGUAGE plpgsql;

select getTax(100);


CREATE or REPLACE FUNCTION getSomeValue(_data json, _key TEXT)
RETURNS TEXT 
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN _data[_key];
END;
$$;

DROP FUNCTION getSomeValue(json,text);

SELECT getSomeValue('{"name":"Allochi", "age":40}','age')

CREATE EXTENSION plpythonu;


CREATE or REPLACE FUNCTION getTax(subtotal REAL) RETURNS REAL AS $$
  return subtotal * 0.07;
$$ LANGUAGE plpythonu;

select getTax(100);

CREATE or REPLACE FUNCTION getJSONValue(_data json, _key TEXT) RETURNS TEXT AS $$
  return type(_data);
$$ LANGUAGE plpythonu;

SELECT getJSONValue('{"name":"Allochi", "age":40}','age')

CREATE OR REPLACE FUNCTION jsel(json_text json, key text)
  RETURNS json
  LANGUAGE plpythonu
  IMMUTABLE
AS $$
if json_text == None or key == None:
    return
import json
j = json.loads(json_text)
tup = key.split(".")
for k in tup:
    if not j.has_key(k):
        return None
    j = j[k]
return j
$$;

SELECT jsel(data, 'more') AS qos FROM "DataStore";
