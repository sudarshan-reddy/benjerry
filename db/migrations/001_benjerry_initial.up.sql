CREATE TABLE ice_cream(
    cursor serial,
    name text PRIMARY KEY,
    image_open text,
    image_closed text,
    story text,
    description text,
    sourcing_values text[],
    ingredients text[],
    allergy_info text,
    dietary_certification text,
    product_id varchar(10)
);
