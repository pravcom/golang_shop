CREATE TABLE IF NOT EXISTS measure_units
(
    id   serial       not null unique primary key,
    code varchar(50) not null unique ,
    name jsonb
);

create table if not exists products
(
    id   serial not null unique primary key,
    name jsonb
);

create table if not exists locations
(
    id      serial not null primary key,
    name    jsonb,
    address varchar(255)
);

create table if not exists orders
(
    id                             serial not null primary key,
    comment                        jsonb,
    source_location_id             serial,
    destination_location_id        serial,
    total_weight_measure_unit_code varchar(50),
    total_volume_measure_unit_code varchar(50),
    foreign key (source_location_id) references locations (id) on update cascade on delete RESTRICT ,
    foreign key (destination_location_id) references locations (id) on update cascade on delete RESTRICT ,
    foreign key (total_weight_measure_unit_code) references measure_units (code) on update cascade ,
    foreign key (total_volume_measure_unit_code) references measure_units (code) on update cascade
);

create table if not exists order_items
(
    id                       serial not null primary key,
    root_id                  serial,
    product_id               serial,
    item_index               integer,
    weight_value             numeric,
    weight_measure_unit_code varchar(50),
    volume_value             numeric,
    volume_measure_unit_code varchar(50),
    foreign key (root_id) references orders (id) on delete CASCADE,
    foreign key (product_id) references products(id) on update cascade on delete RESTRICT ,
    foreign key (weight_measure_unit_code) references measure_units (code) on update cascade ,
    foreign key (volume_measure_unit_code) references measure_units (code) on update cascade
);

begin ;

-- Вставляем единицы измерения
INSERT INTO measure_units (code, name)
VALUES ('kg', '{
  "ru": "килограмм",
  "en": "kilogram"
}'),
       ('g', '{
         "ru": "грамм",
         "en": "gram"
       }'),
       ('t', '{
         "ru": "тонна",
         "en": "ton"
       }'),
       ('l', '{
         "ru": "литр",
         "en": "liter"
       }'),
       ('m3', '{
         "ru": "кубический метр",
         "en": "cubic meter"
       }'),
       ('pcs', '{
         "ru": "штука",
         "en": "piece"
       }');

-- Вставляем продукты
INSERT INTO products (name)
VALUES ('{
  "ru": "Молоко",
  "en": "Milk"
}'),
       ('{
         "ru": "Хлеб",
         "en": "Bread"
       }'),
       ('{
         "ru": "Яблоки",
         "en": "Apples"
       }'),
       ('{
         "ru": "Сахар",
         "en": "Sugar"
       }'),
       ('{
         "ru": "Мясо",
         "en": "Meat"
       }'),
       ('{
         "ru": "Вода",
         "en": "Water"
       }');

-- Вставляем местоположения
INSERT INTO locations (name, address)
VALUES ('{
  "ru": "Склад №1",
  "en": "Warehouse #1"
}', 'ул. Ленина, 10'),
       ('{
         "ru": "Склад №2",
         "en": "Warehouse #2"
       }', 'ул. Гагарина, 25'),
       ('{
         "ru": "Магазин Центральный",
         "en": "Central Store"
       }', 'пр. Мира, 15'),
       ('{
         "ru": "Фабрика",
         "en": "Factory"
       }', 'промзона, сектор 5'),
       ('{
         "ru": "Офис",
         "en": "Office"
       }', 'ул. Садовая, 3');

-- Вставляем заказы
INSERT INTO orders (comment, source_location_id, destination_location_id, total_weight_measure_unit_code,
                    total_volume_measure_unit_code)
VALUES ('{
  "ru": "Срочный заказ",
  "en": "Urgent order"
}', 1, 3, 'kg', 'm3'),
       ('{
         "ru": "Регулярная поставка",
         "en": "Regular delivery"
       }', 4, 2, 't', 'm3'),
       ('{
         "ru": "Междугородняя доставка",
         "en": "Intercity delivery"
       }', 2, 5, 'kg', 'l'),
       ('{
         "ru": "Заказ для магазина",
         "en": "Store order"
       }', 3, 1, 'g', 'm3');

-- Вставляем элементы заказов
INSERT INTO order_items (root_id, product_id, item_index, weight_value, weight_measure_unit_code, volume_value,
                         volume_measure_unit_code)
VALUES (1, 1, 1, 10.5, 'kg', 0.01, 'm3'),
       (1, 2, 2, 2.0, 'kg', 0.002, 'm3'),
       (1, 3, 3, 5.0, 'kg', 0.005, 'm3'),
       (2, 4, 1, 500.0, 'kg', 0.3, 'm3'),
       (2, 5, 2, 200.0, 'kg', 0.15, 'm3'),
       (3, 6, 1, 15.0, 'kg', 0.015, 'm3'),
       (3, 1, 2, 8.0, 'kg', 0.008, 'm3'),
       (4, 2, 1, 1000.0, 'g', 0.001, 'm3'),
       (4, 3, 2, 1500.0, 'g', 0.0015, 'm3'),
       (4, 4, 3, 800.0, 'g', 0.0008, 'm3');

commit