import argparse
import random
import string
import re
import os


def parse_size(size_str: str) -> int:
    """
    Парсит строку размера вида '10kb', '1mb', '512b' и возвращает размер в байтах.
    """
    match = re.fullmatch(r'(?i)(\d+(?:\.\d+)?)(b|kb|mb|gb)', size_str.strip())
    if not match:
        raise ValueError(f"Неверный формат размера: {size_str}")

    size, unit = match.groups()
    size = float(size)
    unit = unit.lower()
    unit_multipliers = {
        'b': 1,
        'kb': 1024,
        'mb': 1024 ** 2,
        'gb': 1024 ** 3,
    }

    return int(size * unit_multipliers[unit])


def generate_string_to_file(byte_size: int, filepath: str) -> None:
    """
    Генерирует строку из латинских букв заданного размера (в байтах) и записывает в файл.
    """
    chars = string.ascii_letters
    chunk_size = 1024 * 1024  # 1MB

    os.makedirs(os.path.dirname(filepath), exist_ok=True)

    with open(filepath, 'w', encoding='utf-8') as f:
        while byte_size > 0:
            this_chunk = min(byte_size, chunk_size)
            f.write(''.join(random.choices(chars, k=this_chunk)))
            byte_size -= this_chunk

    print(f"File '{filepath}' generated ({os.path.getsize(filepath)} bytes)")


def main():
    parser = argparse.ArgumentParser(description="Генерация строк заданного размера и сохранение в файлы.")
    parser.add_argument("sizes", nargs="+", help="Размеры в формате 10kb, 2mb, 512b и т.д.")
    parser.add_argument("--output-dir", "-o", default=".", help="Каталог для сохранения файлов (по умолчанию текущий).")

    args = parser.parse_args()

    for idx, size_str in enumerate(args.sizes, 1):
        try:
            byte_size = parse_size(size_str)
            filename = f"output_{idx}.txt"
            filepath = os.path.join(args.output_dir, filename)
            generate_string_to_file(byte_size, filepath)
        except ValueError as e:
            print(f"⚠️ Ошибка: {e}")


if __name__ == "__main__":
    main()
