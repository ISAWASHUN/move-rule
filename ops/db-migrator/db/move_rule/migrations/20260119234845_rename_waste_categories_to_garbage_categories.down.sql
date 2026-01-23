-- 外部キー制約を削除
ALTER TABLE garbage_items DROP FOREIGN KEY garbage_items_ibfk_2;

-- カラム名を元に戻す
ALTER TABLE garbage_items CHANGE garbage_category_id waste_category_id INT NOT NULL COMMENT '分別区分ID';

-- テーブル名を元に戻す
RENAME TABLE garbage_categories TO waste_categories;

-- 外部キー制約を再作成
ALTER TABLE garbage_items ADD CONSTRAINT garbage_items_ibfk_2 FOREIGN KEY (waste_category_id) REFERENCES waste_categories(id) ON DELETE RESTRICT;
