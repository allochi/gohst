-- create extension json_enhancements;

Delete FROM json_contacts;
TRUNCATE TABLE json_contacts RESTART IDENTITY;

DROP VIEW json_contacts_view;
CREATE VIEW json_contacts_view AS
SELECT 
contacts.*, 
titles."name" as title,
job_titles."name" as job_title,
countries."printable_name" as country,
cities."name" as city,
organizations."name" as organization,
departments."name" as department,
Array[contacts.email_1, contacts.email_2] as emails,
contacts."isOrganization"::INT::BOOL as is_organization,
contacts."focal_point"::INT::BOOL as is_focal_point
FROM contacts 
LEFT JOIN titles ON titles."id" = contacts.title_id
LEFT JOIN job_titles ON job_titles."id" = contacts.job_title_id
LEFT JOIN countries ON countries."id" = contacts.country_id
LEFT JOIN cities ON cities."id" = contacts.city_id
LEFT JOIN organizations ON organizations."id" = contacts.organization_id
LEFT JOIN departments ON departments."id" = contacts.department_id;

INSERT INTO json_contacts 
SELECT id as Id, row_to_json(json_contacts_view) as Data, created_at as CreatedAt, updated_at as UpdateAt FROM json_contacts_view;
