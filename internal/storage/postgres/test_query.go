package postgres

import (
	"context"
	"fmt"
	"github.com/2pizzzza/plumbing/internal/storage"
)

func (db *DB) CreateStarter(ctx context.Context) error {
	const op = "postgres.CreateStarter"

	var exists bool
	checkLanguageQuery := `SELECT EXISTS(SELECT 1 FROM Language WHERE code IN ('ru', 'kgz', 'en'))`
	err := db.Pool.QueryRow(ctx, checkLanguageQuery).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: failed to check if languages already exist: %w", op, err)
	}
	if exists {
		return storage.ErrCategoryNotFound
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO Language (code, name) VALUES 
		('ru', 'Русский'),
		('kgz', 'Кырғызча'),
		('en', 'English')`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert languages: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO Category (id) VALUES (1), (2), (3)`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert categories: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO CategoryTranslation (category_id, language_code, name) VALUES
		(1, 'ru', 'Техника'),
		(1, 'kgz', 'Техника'),
		(1, 'en', 'Electronics'),
		(2, 'ru', 'Одежда'),
		(2, 'kgz', 'Кийгени кийим'),
		(2, 'en', 'Clothing'),
		(3, 'ru', 'Обувь'),
		(3, 'kgz', 'Атайын бут кийим'),
		(3, 'en', 'Footwear')`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert category translations: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO Collection (price, isProducer, isPainted, isPopular, isNew) 
		VALUES 
		(500.00, true, true, true, false), 
		(1000.00, false, false, false, true),
		(1200.00, true, false, false, false),
		(1500.00, false, true, true, true)`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert collections: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO CollectionTranslation (collection_id, language_code, name, description) VALUES
		(1, 'ru', 'Смарт телевизор', 'Современный телевизор с высокой четкостью изображения и умными функциями'),
		(1, 'kgz', 'Смарт телевизор', 'Жогорку сапаттагы сүрөттү жана акылдуу функцияларды камтыган заманбап телевизор'),
		(1, 'en', 'Smart TV', 'A modern TV with high-definition display and smart features'),
		(2, 'ru', 'Летняя коллекция', 'Легкая и удобная одежда для жаркой погоды'),
		(2, 'kgz', 'Жайкы коллекция', 'Ысык аба ырайы үчүн жеңил жана ыңгайлуу кийимдер'),
		(2, 'en', 'Summer Collection', 'Light and comfortable clothing for hot weather'),
		(3, 'ru', 'Осенние новинки', 'Теплая одежда для осени'),
		(3, 'kgz', 'Күзгү жаңы коллекция', 'Күзгө ылайыктуу жылуу кийимдер'),
		(3, 'en', 'Autumn New Arrivals', 'Warm clothing for autumn'),
		(4, 'ru', 'Весенняя обувь', 'Обувь для весны и лета'),
		(4, 'kgz', 'Жазкы бут кийим', 'Жазга жана жайга ылайыктуу бут кийимдер'),
		(4, 'en', 'Spring Footwear', 'Footwear for spring and summer')`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert collection translations: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO Item (category_id, collection_id, size, price, isProducer, isPainted, isPopular, isNew) 
		VALUES 
		(1, 1, '50 inch', 500.00, true, true, true, false),
		(2, 2, 'M', 1000.00, false, false, false, true),
		(3, 4, '42', 1500.00, false, true, true, true),
		(1, 3, '55 inch', 1200.00, true, false, false, false),
		(2, 2, 'L', 800.00, false, false, true, true),
		(1, 1, '65 inch', 700.00, true, true, false, true),
		(3, 4, '40', 1600.00, true, true, true, false)`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert items: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO ItemTranslation (item_id, language_code, name, description) VALUES
		(1, 'ru', '4K Смарт ТВ', 'Телевизор с поддержкой 4K и интернет-возможностями'),
		(1, 'kgz', '4K Смарт ТВ', '4K колдонуучу жана интернет мүмкүнчүлүктөрү бар телевизор'),
		(1, 'en', '4K Smart TV', 'TV with 4K support and internet features'),
		(2, 'ru', 'Летнее платье', 'Легкое платье для жаркой погоды'),
		(2, 'kgz', 'Жайкы көйнөк', 'Ысык күнгө ылайыктуу жеңил көйнөк'),
		(2, 'en', 'Summer Dress', 'A light dress for hot weather'),
		(3, 'ru', 'Кожаные туфли', 'Легкие туфли из натуральной кожи для осени'),
		(3, 'kgz', 'Жайкы бут кийим', 'Күзгө ылайыктуу табигый териден жасалган жеңил бут кийимдер'),
		(3, 'en', 'Leather Shoes', 'Light leather shoes for autumn')`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert item translations: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO Photo (url, isMain) VALUES 
		('media/images/users.png', true),
		('media/images/users.png', false),
		('media/images/users.png', true),
		('media/images/users.png', false)`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert photos: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO CollectionPhoto (collection_id, photo_id) VALUES
		(1, 1), 
		(2, 2),
		(3, 3),
		(4, 4)`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert collection photos: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO ItemPhoto (item_id, photo_id) VALUES
		(1, 1),
		(2, 2),
		(3, 3),
		(4, 4),
		(5, 1)`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert item photos: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO Color (hash_color) VALUES 
		('#FF5733'), 
		('#33FF57'),
		('#C70039'), 
		('#900C3F')`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert colors: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO CollectionColor (collection_id, color_id) VALUES
		(1, 1), 
		(2, 2),
		(3, 3),
		(4, 4)`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert collection colors: %w", op, err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO ItemColor (item_id, color_id) VALUES
		(1, 1),
		(2, 2),
		(3, 3),
		(4, 4),
		(5, 1)`)
	if err != nil {
		return fmt.Errorf("%s: failed to insert item colors: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}
