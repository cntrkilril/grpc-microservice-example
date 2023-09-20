CREATE TYPE CANCELLED_BY AS ENUM ('admin', 'master', 'client');

CREATE TABLE IF NOT EXISTS appointment
(
    id              SERIAL PRIMARY KEY,
    master_id       INT  NOT NULL,
    client_id       INT  NOT NULL,
    service_id      INT  NOT NULL,
    start_time      TIME NOT NULL,
    end_time        TIME NOT NULL,
    date            DATE NOT NULL,
    is_confirmed    BOOL NOT NULL,
    cancelled_at    TIMESTAMP,
    cancel_reason   TEXT,
    cancelled_by    CANCELLED_BY,
    cancelled_by_id INT
);

CREATE TABLE IF NOT EXISTS time_cell
(
    id             SERIAL PRIMARY KEY,
    master_id      INT  NOT NULL,
    date           DATE NOT NULL,
    start_time     TIME NOT NULL,
    end_time       TIME NOT NULL,
    is_free        BOOL NOT NULL,
    appointment_id INT,
    FOREIGN KEY (appointment_id) REFERENCES appointment (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS work_time
(
    id         SERIAL PRIMARY KEY,
    master_id  INT  NOT NULL,
    start_time TIME NOT NULL,
    end_time   TIME NOT NULL,
    date       DATE NOT NULL
);