-- This is a test migration for the dev environment 
CREATE TABLE ingenious (
    orderId VARCHAR(255) NOT NULL,
    barcode VARCHAR(255) NOT NULL,
    usi VARCHAR(255) NOT NULL,
    patient VARCHAR(255) NOT NULL,
    birthdate DATE NOT NULL
);

-- Insert dummy sample data
INSERT INTO ingenious (orderId, barcode, usi, patient, birthdate) VALUES
('123456789', '123456789', '123456789', 'John Doe', '1990-01-01'),
('987654321', '987654321', '987654321', 'Jane Doe', '1990-01-01');
