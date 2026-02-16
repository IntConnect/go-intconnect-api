INSERT INTO alarm_logs(parameter_id, acknowledged_by, value, type, note)
VALUES (1, null, 1, 'Open', '')
SELECT
   table_name,
   column_name,
   data_type
FROM
   information_schema.columns
WHERE
   table_name = 'users';