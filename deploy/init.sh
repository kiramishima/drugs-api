#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE DATABASE credits;

  CREATE TABLE IF NOT EXISTS credits.credit_assigns(
       id BIGSERIAL NOT NULL PRIMARY KEY,
       invest INT NOT NULL DEFAULT 0,
       credit_300 INT NOT NULL DEFAULT 0,
       credit_500 INT NOT NULL DEFAULT 0,
       credit_700 INT NOT NULL DEFAULT 0,
       status SMALLINT NOT NULL DEFAULT 0,
       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP,
       deleted_at TIMESTAMP
  );

  CREATE VIEW credits.statistics AS
  WITH gral_credits AS (
      SELECT COUNT(*) total, SUM(invest) total_inv FROM credits.credit_assigns
  ), success_credits AS (
      SELECT COUNT(*) total_sucess, SUM(invest) total_success_inv FROM credits.credit_assigns
      WHERE status = 1
  ), fail_credits AS (
      SELECT COUNT(*) total_fails, SUM(invest) total_fail_inv FROM credits.credit_assigns
      WHERE status = 0
  )
  SELECT gral_credits.total,
         gral_credits.total_inv,
         total_sucess,
         total_fails,
         success_credits.total_success_inv,
         ROUND(((success_credits.total_success_inv / gral_credits.total_inv::float)*100)::numeric, 2) avg_total_success_inv,
         ROUND(((total_fail_inv/total_inv::float)*100)::numeric, 2) avg_total_fail_inv
  FROM gral_credits, success_credits, fail_credits;
EOSQL