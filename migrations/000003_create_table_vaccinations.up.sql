CREATE TABLE IF NOT EXISTS vaccinations(
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    drug_id INTEGER NOT NULL,
    dose SMALLINT NOT NULL,
    applied_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE (name, drug_id, applied_at),
    CONSTRAINT fk_drugs
        FOREIGN KEY (drug_id)
            REFERENCES drugs(id)
);