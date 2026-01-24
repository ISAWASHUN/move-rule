CREATE TABLE garbage_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    municipality_id INT NOT NULL COMMENT '地方公共団体ID',
    garbage_category_id INT NOT NULL COMMENT '分別区分ID',
    area_name VARCHAR(255) COMMENT '地区名',
    item_name VARCHAR(255) NOT NULL COMMENT 'ゴミの品目',
    notes TEXT COMMENT '注意点',
    remarks TEXT COMMENT '備考',
    bulk_garbage_fee INT COMMENT '粗大ごみ回収料金',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード更新日時',
    FOREIGN KEY (municipality_id) REFERENCES municipalities(id) ON DELETE CASCADE,
    FOREIGN KEY (garbage_category_id) REFERENCES garbage_categories(id) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ゴミの品目';
