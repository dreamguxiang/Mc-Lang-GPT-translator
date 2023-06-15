package langparser

import (
	"bufio"
	"os"
	"strings"
)

type LangEntry struct {
	Key   string
	Value string
}

func ParseLangFile(filePath string) ([]LangEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []LangEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			entry := LangEntry{
				Key:   strings.TrimSpace(parts[0]),
				Value: strings.TrimSpace(parts[1]),
			}
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func ReplaceLangEntry(oldEntries []LangEntry, newEntries []LangEntry) []LangEntry {
	for i, oldEntry := range oldEntries {
		for _, newEntry := range newEntries {
			if oldEntry.Key == newEntry.Key {
				oldEntries[i].Value = newEntry.Value
			}
		}
	}
	return oldEntries
}

func ReadLangByString(data string) ([]LangEntry, error) {
	var entries []LangEntry
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			entry := LangEntry{
				Key:   strings.TrimSpace(parts[0]),
				Value: strings.TrimSpace(parts[1]),
			}
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil

}

func WriteLangFile(filePath string, entries []LangEntry) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, entry := range entries {
		line := entry.Key + "=" + entry.Value + "\n"
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
