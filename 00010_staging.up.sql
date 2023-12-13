CREATE TABLE IF NOT EXISTS ingenious
(
    orderid character varying,
    barcode character varying NOT NULL,
    usi character varying NOT NULL,
    patient character varying,
    specimen character varying,
    birthdate date
);

-- insert dummy data
INSERT INTO ingenious (orderid, barcode, usi, patient, birthdate, specimen) VALUES ('1233031845', '123303184503', 'CHLAPPA', 'AMIRI,ANNA', '2020-01-01', 'BAL');
INSERT INTO ingenious (orderid, barcode, usi, patient, birthdate, specimen) VALUES ('1233031851', '123303185103', 'CHLAPPA', 'T,TEST', '2020-01-01', 'BAL');
INSERT INTO ingenious (orderid, barcode, usi, patient, birthdate, specimen) VALUES ('1233031852', '123303185203', 'CHLAPPA', 'T,ADSADAS', '2020-01-01', 'BAL');
INSERT INTO ingenious (orderid, barcode, usi, patient, birthdate, specimen) VALUES ('1233031853', '123303185303', 'CHLAPPA', 'T,Kalle', '2020-01-01', 'BAL');

INSERT INTO samples (sample_id, full_name, birthdate, comment, manual) VALUES ('123303184503', 'AMIRI,ANNA', '2020-01-01', 'comment testing', true);