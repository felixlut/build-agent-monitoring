CREATE TABLE webhooks_workflow_job (
    workflow_run_id BIGINT NOT NULL,
    workflow_job_status VARCHAR(32) NOT NULL,
    workflow_job_id BIGINT NOT NULL
    -- conclusion VARCHAR(32),
    -- started_at TIMESTAMP,
    -- concluded_at TIMESTAMP,
    -- labels json NOT NULL,
    -- runner_id INT,
    -- runner_groud_id INT,
    -- runner_name VARCHAR(64),
    -- runner_group_name VARCHAR(64),
    -- repository_id BIGINT NOT NULL,
    -- repository_NAME VARCHAR(64) NOT NULL,
    -- organization_id BIGINT NOT NULL,
    -- organization_NAME VARCHAR(64) NOT NULL,
    -- sender_id BIGINT NOT NULL,
    -- sender_NAME VARCHAR(64) NOT NULL
);
