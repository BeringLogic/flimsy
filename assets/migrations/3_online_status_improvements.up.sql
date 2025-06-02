ALTER TABLE config ADD online_status_timeout INTEGER;
UPDATE config SET online_status_timeout = 10;

ALTER TABLE item ADD skip_certificate_verification INTEGER;
UPDATE item SET skip_certificate_verification = 0;

ALTER TABLE item ADD check_url TEXT;
UPDATE item SET check_url = url;
