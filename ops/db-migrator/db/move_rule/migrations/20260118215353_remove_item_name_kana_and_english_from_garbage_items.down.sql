ALTER TABLE garbage_items ADD COLUMN item_name_kana VARCHAR(255) COMMENT 'ゴミの品目_カナ' AFTER item_name;
ALTER TABLE garbage_items ADD COLUMN item_name_english VARCHAR(255) COMMENT 'ゴミの品目_英字' AFTER item_name_kana;
