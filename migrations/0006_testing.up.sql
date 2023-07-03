CREATE TABLE IF NOT EXISTS ingenious(
    orderId VARCHAR (255) NOT NULL,
    barcode VARCHAR (255) NOT NULL,
    usi VARCHAR (255) NOT NULL,
    patient VARCHAR (255) NOT NULL,
    birthdate DATE NOT NULL
);

-- Dummy samples for testing
INSERT INTO ingenious (orderId, barcode, usi, patient, birthdate) VALUES ('1232088718', '123208871803', 'MYPPA', 'John Doe', '1990-01-01');
INSERT INTO ingenious (orderId, barcode, usi, patient, birthdate) VALUES ('1232088719', '123208871903', 'MYPPA', 'Jane Doe', '1990-01-01');
INSERT INTO ingenious (orderId, barcode, usi, patient, birthdate) VALUES ('1232088720', '123208872003', 'CHLAPPA', 'John2 Doe', '1990-01-01');
INSERT INTO ingenious (orderId, barcode, usi, patient, birthdate) VALUES ('1232088721', '123208872103', 'MYPPA', 'Jane1 Doe', '1990-01-01');
INSERT INTO ingenious (orderId, barcode, usi, patient, birthdate) VALUES ('1232088722', '123208872203', 'CHLAPPA', 'John2 Doe', '1990-01-01');