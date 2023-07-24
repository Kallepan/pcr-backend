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
('123456789', '123456789', 'AFAG', 'John Doe', '1990-01-01'),
('987654321', '987654321', 'KLASDN', 'Jane Doe', '1990-01-01'),
('531313115', '531313115', 'ASDASD', 'John Smith', '1990-01-01'),
('158777887', '158777887', 'ASDASD', 'Jane Smith', '1990-01-01'),
('455454454', '455454454', 'ASDASD', 'KASADAS', '1990-01-01');
