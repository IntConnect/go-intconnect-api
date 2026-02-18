INSERT INTO alarm_logs(parameter_id, acknowledged_by, value, type, note)
VALUES (1, null, 1, 'Open', '')
SELECT table_name,
       column_name,
       data_type
FROM information_schema.columns
WHERE table_name = 'users';
SELECT * FROM telemetries WHERE timestamp > '2026-02-18 21:00:00' ORDER BY id DESC;
SELECT bucket,
       parameter_id,
       last_value
FROM (SELECT time_bucket_gapfill('1 minutes'::interval, timestamp) AS bucket,
             parameter_id,
             last(value, timestamp)                                AS last_value
      FROM telemetries
      WHERE parameter_id IN (3, 4, 5)
        AND timestamp BETWEEN '2026-02-18 21:00:00' AND '2026-02-18 21:30:00'
      GROUP BY bucket, parameter_id) q
ORDER BY bucket;