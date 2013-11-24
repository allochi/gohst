-- More practical is just string processing
CREATE OR REPLACE FUNCTION _array(_j json, _key text) RETURNS text[] as $$
    SELECT concat('{',btrim(_j->>_key,'[]'),'}')::text[]
$$ LANGUAGE SQL IMMUTABLE;

-- DROP FUNCTION _array(_j json, _key text);
-- CREATE OR REPLACE FUNCTION _array(_j json, _key text) RETURNS text[] as $$
--     SELECT array_agg(btrim(x.elem::text, '"')) from json_array_elements((_j->>_key)::json) as x(elem)
-- $$ LANGUAGE SQL IMMUTABLE;

DROP INDEX _array_categories_idx;
CREATE INDEX _array_categories_idx ON json_contacts USING GIN (_array(data,'categories'));

-- Total runtime: 37.412 ms no index (9.3.1 => Total runtime: 17.785 ms)
-- Total runtime: 0.449 ms index (100x) (9.3.1 => Total runtime: 0.238 ms => 7x)
EXPLAIN ANALYZE SELECT _array(data,'categories'), data->>'categories' FROM json_contacts
WHERE '{Potential Donors}' <@ _array(data,'categories');

SELECT _array(data,'interests') FROM json_contacts;

select array_agg(z) from 
(SELECT json_array_elements(_temp.i) FROM
(SELECT data->'interests' as i FROM json_contacts WHERE id = 1) as _temp) as z;


-- Date and Date Array
DROP FUNCTION IF EXISTS _date(_j json, _key text);
CREATE OR REPLACE FUNCTION _date(_j json, _key text) RETURNS TIMESTAMP as $$
    SELECT (_j->>_key)::TIMESTAMP
$$ LANGUAGE SQL IMMUTABLE;

DROP FUNCTION IF EXISTS _date_array(_j json, _key text);
CREATE OR REPLACE FUNCTION _date_array(_j json, _key text) RETURNS TIMESTAMP[] as $$
    SELECT array_agg(x.elem::TIMESTAMP) from json_array_elements((_j->>_key)::json) as x(elem)
$$ LANGUAGE SQL IMMUTABLE;

