-- Explain
EXPLAIN ANALYZE SELECT * FROM contacts WHERE ID = 1;

-- Explain
EXPLAIN ANALYZE SELECT * FROM contacts WHERE first_name = 'Ali' AND id = 400;

-- Index
CREATE INDEX first_name_idx ON contacts (first_name);

-- Index
CREATE INDEX id_first_name_idx ON contacts (id,first_name);

-- Function
CREATE or REPLACE FUNCTION merge_with_id_field(id REAL, field TEXT) RETURNS TEXT AS $$
BEGIN
  RETURN id || field;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Test
SELECT merge_with_id_field(400,'Ali');

-- Index
CREATE INDEX merge_id_first_name_idx ON contacts (merge_with_id_field(id,first_name));

-- Explain
EXPLAIN ANALYZE SELECT * FROM contacts WHERE merge_with_id_field(id,first_name) = '400Ali';

-- Function
DROP FUNCTION getCountry(_id REAL);
CREATE or REPLACE FUNCTION getCountry(_id REAL) RETURNS TEXT AS $$
  SELECT name FROM countries WHERE id = _id;
$$ LANGUAGE SQL IMMUTABLE;

SELECT getCountry(4);

-- Index
DROP INDEX country_idx;
CREATE INDEX country_idx ON contacts (getCountry(country_id));

EXPLAIN ANALYZE SELECT * FROM contacts WHERE getCountry(country_id) = 'CANADA';

SELECT countries.name FROM contacts left join countries on contacts.country_id = countries.id;

----------------------------------------
-- JSONify contacts
----------------------------------------
create extension json_enhancements;
INSERT INTO json_contacts 
SELECT id, created_at, updated_at,row_to_json(contacts) as data FROM contacts;


-- Get country_id
SELECT (data->>'country_id')::int FROM json_contacts;

-- Function
DROP FUNCTION getJSONCountry(data json);
CREATE or REPLACE FUNCTION getJSONCountry(data json) RETURNS TEXT AS $$
  SELECT name FROM countries WHERE id = (data->>'country_id')::int;
$$ LANGUAGE SQL IMMUTABLE;

-- Index
DROP INDEX json_country_idx;
CREATE INDEX json_country_idx ON json_contacts (getJSONCountry(data));


EXPLAIN ANALYZE SELECT * FROM json_contacts WHERE getJSONCountry(data) = 'CANADA';

-- Make 'name' from contacts without indexes, costs AVG. 0.606ms
EXPLAIN ANALYZE SELECT (first_name || ' ' || last_name) as name FROM contacts;

-- Make 'name' from json_contacts without indexes, costs AVG. 39.392ms
EXPLAIN ANALYZE SELECT concat_ws(' ',data->>'first_name', data->>'last_name') as name FROM json_contacts;

EXPLAIN (format json) SELECT concat_ws(' ',data->>'first_name', data->>'last_name') as name FROM json_contacts;

----------------------------------------
-- SELECT the Index!
----------------------------------------
CREATE INDEX first_name_idx ON contacts (first_name);
SELECT first_name_idx FROM contacts FETCH FIRST 100 ROWS ONLY;
