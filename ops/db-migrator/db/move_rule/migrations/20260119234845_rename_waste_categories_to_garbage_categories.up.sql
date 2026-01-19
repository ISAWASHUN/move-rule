-- waste_categoriesテーブルをgarbage_categoriesにリネーム
RENAME TABLE waste_categories TO garbage_categories;

-- garbage_itemsテーブルのカラム名を変更
ALTER TABLE garbage_items CHANGE waste_category_id garbage_category_id INT NOT NULL COMMENT '分別区分ID';

-- 外部キー制約を削除して再作成
ALTER TABLE garbage_items DROP FOREIGN KEY garbage_items_ibfk_2;
ALTER TABLE garbage_items ADD CONSTRAINT garbage_items_ibfk_2 FOREIGN KEY (garbage_category_id) REFERENCES garbage_categories(id) ON DELETE RESTRICT;
