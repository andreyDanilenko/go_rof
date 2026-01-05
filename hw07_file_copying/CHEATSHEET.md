# Шпаргалка: Основные концепции и паттерны

Быстрая справка по ключевым концепциям для реализации Copy.

## 📚 Основные функции и методы

### Работа с файлами

```go
// Открыть файл для чтения
file, err := os.Open("path/to/file")
defer file.Close()  // ВАЖНО: всегда закрывать!

// Создать/перезаписать файл для записи
file, err := os.Create("path/to/file")
defer file.Close()

// Получить информацию о файле
info, err := file.Stat()
size := info.Size()  // размер в байтах (int64)

// Проверить, что это обычный файл
isRegular := info.Mode().IsRegular()
// или
if info.Mode()&os.ModeType == 0 {
    // это обычный файл
}
```

### Перемещение по файлу

```go
import "io"

// Переместиться на позицию offset от начала файла
newPos, err := file.Seek(offset, io.SeekStart)

// Текущая позиция (без перемещения)
pos, err := file.Seek(0, io.SeekCurrent)

// Константы для Seek:
io.SeekStart   // от начала файла
io.SeekCurrent // от текущей позиции
io.SeekEnd     // от конца файла
```

### Копирование данных

```go
import "io"

// Скопировать максимум n байт из src в dst
written, err := io.CopyN(dst, src, n)
// Возвращает количество скопированных байт и ошибку

// Пример:
written, err := io.CopyN(outputFile, inputFile, 1000)
if err != nil {
    return err
}
// written может быть меньше 1000, если в src осталось меньше данных
```

### Прогресс-бар (с библиотекой)

```go
import "github.com/cheggaaa/pb/v3"

// Создать прогресс-бар
bar := pb.StartNew(totalBytes)

// Обернуть Reader для отслеживания прогресса
reader := bar.NewProxyReader(file)

// Использовать обернутый reader для копирования
io.CopyN(outputFile, reader, copySize)

// Завершить прогресс-бар
bar.Finish()
```

---

## 🔢 Вычисление размера копирования

```go
// Шаг 1: Вычислить оставшийся размер
remaining := fileSize - offset

// Шаг 2: Определить сколько копировать
var copySize int64
if limit == 0 {
    // Копируем все, что осталось
    copySize = remaining
} else {
    // Копируем минимум из limit и остатка
    if limit < remaining {
        copySize = limit
    } else {
        copySize = remaining
    }
}

// Или короче с использованием math.Min:
import "math"
copySize := remaining
if limit > 0 && limit < remaining {
    copySize = limit
}
```

---

## ✅ Паттерны проверок и валидации

### Паттерн 1: Открыть и проверить

```go
file, err := os.Open(path)
if err != nil {
    return err
}
defer file.Close()  // всегда откладываем закрытие
```

### Паттерн 2: Получить информацию и проверить

```go
info, err := file.Stat()
if err != nil {
    return err
}

size := info.Size()
if !info.Mode().IsRegular() {
    return ErrUnsupportedFile
}
```

### Паттерн 3: Валидация входных параметров

```go
if offset < 0 {
    return errors.New("offset cannot be negative")
}

if offset >= fileSize {
    return ErrOffsetExceedsFileSize
}
```

### Паттерн 4: Выполнить операцию и проверить результат

```go
newPos, err := file.Seek(offset, io.SeekStart)
if err != nil {
    return err
}

if newPos != offset {
    return fmt.Errorf("seek failed: expected %d, got %d", offset, newPos)
}
```

---

## 🔄 Правильный порядок операций

```go
func Copy(fromPath, toPath string, offset, limit int64) error {
    // 1. Открыть исходный файл
    src, err := os.Open(fromPath)
    if err != nil {
        return err
    }
    defer src.Close()
    
    // 2. Получить информацию о файле
    info, err := src.Stat()
    if err != nil {
        return err
    }
    
    // 3. Проверить тип файла
    if !info.Mode().IsRegular() {
        return ErrUnsupportedFile
    }
    
    // 4. Проверить offset
    fileSize := info.Size()
    if offset >= fileSize {
        return ErrOffsetExceedsFileSize
    }
    
    // 5. Переместиться на offset
    _, err = src.Seek(offset, io.SeekStart)
    if err != nil {
        return err
    }
    
    // 6. Вычислить размер копирования
    remaining := fileSize - offset
    copySize := remaining
    if limit > 0 && limit < remaining {
        copySize = limit
    }
    
    // 7. Создать файл назначения
    dst, err := os.Create(toPath)
    if err != nil {
        return err
    }
    defer dst.Close()
    
    // 8. Скопировать данные (с прогресс-баром или без)
    // ...
    
    return nil
}
```

---

## 📊 Структура для прогресс-бара (свой вариант)

```go
type progressWriter struct {
    writer  io.Writer
    total   int64
    written int64
}

func (pw *progressWriter) Write(p []byte) (int, error) {
    n, err := pw.writer.Write(p)
    pw.written += int64(n)
    
    // Вычислить и вывести процент
    percent := (pw.written * 100) / pw.total
    fmt.Printf("\rПрогресс: %d%%", percent)
    
    return n, err
}

// Использование:
pw := &progressWriter{
    writer: dstFile,
    total:  copySize,
}
io.CopyN(pw, srcFile, copySize)
fmt.Println() // новая строка после завершения
```

---

## 🧪 Паттерны для тестов

### Создание тестовых файлов

```go
func TestExample(t *testing.T) {
    // Создать временный исходный файл
    src, err := os.CreateTemp("", "test-src-*.txt")
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(src.Name())
    defer src.Close()
    
    // Записать тестовые данные
    testData := []byte("Hello, World!")
    src.Write(testData)
    src.Close()
    
    // Создать временный файл для результата
    dst, err := os.CreateTemp("", "test-dst-*.txt")
    if err != nil {
        t.Fatal(err)
    }
    dst.Close()
    defer os.Remove(dst.Name())
    
    // Выполнить копирование
    err = Copy(src.Name(), dst.Name(), 0, 0)
    if err != nil {
        t.Fatal(err)
    }
    
    // Проверить результат
    result, err := os.ReadFile(dst.Name())
    if err != nil {
        t.Fatal(err)
    }
    
    if !bytes.Equal(result, testData) {
        t.Errorf("expected %q, got %q", testData, result)
    }
}
```

### Проверка ошибок

```go
err := Copy("nonexistent.txt", "out.txt", 0, 0)
if err == nil {
    t.Error("expected error for nonexistent file")
}

err = Copy("test.txt", "out.txt", 10000, 0)
if !errors.Is(err, ErrOffsetExceedsFileSize) {
    t.Errorf("expected ErrOffsetExceedsFileSize, got %v", err)
}
```

---

## 💡 Частые ошибки и как их избежать

### ❌ Ошибка 1: Забыли закрыть файл

```go
// ПЛОХО:
file, _ := os.Open("file.txt")
// файл не закрыт - утечка ресурсов!

// ХОРОШО:
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close()  // закроется автоматически
```

### ❌ Ошибка 2: Не проверили offset

```go
// ПЛОХО:
src.Seek(offset, io.SeekStart)  // может выйти за границы!

// ХОРОШО:
if offset >= fileSize {
    return ErrOffsetExceedsFileSize
}
src.Seek(offset, io.SeekStart)
```

### ❌ Ошибка 3: Неправильное вычисление copySize

```go
// ПЛОХО:
copySize := limit  // если limit=0 или limit > остаток, это неправильно!

// ХОРОШО:
remaining := fileSize - offset
copySize := remaining
if limit > 0 && limit < remaining {
    copySize = limit
}
```

### ❌ Ошибка 4: Не проверили результат Seek

```go
// ПЛОХО:
src.Seek(offset, io.SeekStart)  // игнорируем ошибку!

// ХОРОШО:
newPos, err := src.Seek(offset, io.SeekStart)
if err != nil {
    return err
}
```

---

## 🎯 Ключевые моменты для запоминания

1. **Всегда используйте `defer` для закрытия файлов**
   - Гарантирует закрытие даже при ошибке
   - Ставьте `defer` сразу после проверки ошибки открытия

2. **Порядок операций важен**
   - Сначала проверки (размер, тип файла)
   - Потом операции (Seek, Copy)

3. **limit = 0 означает "все"**
   - Это удобное соглашение
   - Обрабатывайте отдельно

4. **Проверяйте все ошибки**
   - Не игнорируйте возвращаемые ошибки
   - Особенно от Seek и CopyN

5. **Прогресс-бар оборачивает Reader**
   - Не Writer, а именно Reader
   - Прогресс отслеживается при чтении

---

## 📖 Полезные ссылки

- [Go документация: os](https://pkg.go.dev/os)
- [Go документация: io](https://pkg.go.dev/io)
- [Go документация: io/fs](https://pkg.go.dev/io/fs)
- [pb библиотека](https://github.com/cheggaaa/pb)

---

**Используйте эту шпаргалку как справочник при написании кода!**


