package other

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

const maxFileSize int64 = 1e+8

var ignoreRules = []string{
	"xtensionxtooltip2",
	"Новый владелец канала",
	"изменяет пароль",
	"Creature Removed",
	"не был поставлен менестрелем",
	"Смена канала:",
	"You set all speeds of",
	"You are summoning",
	"Вы вне роли, оставив сообщение",
	"Вы вышли из режима",
	"Game Object",
	"Рысканье",
	"Set",
	"Вы заслужили новое звание",
	"SpawnTime:",
	"PhaseGroup",
	"PhaseId",
	"Orientation",
	"Selected object:",
	"Syntax:",
	"GUID: ",
	"You set the size",
	"PhaseShift:",
	"GroundZ:",
	"ZoneX:",
	"grid[",
	"Map: ",
	"Fly Mode",
	"You are outdoors",
	"UiWorldMapAreaSwaps",
	"VisibleMapIds",
	"PersonalGuid",
	"X: ",
	"Accepting Whisper: ",
	"Darkmoon 905",
	"Вы приобрели новую способность: ",
	"выходит из игрового мира.",
	"Вы присоединились к рейдовой группе",
	"Вы покидаете группу.",
	"Установленный режим сложности подземелья:",
	"Appearing at ",
	"You can only summon a player to your instance",
	"Entry: ",
	"Position: ",
	"Descrption: ",
	"Description: ",
	"Type: ",
	"Name: ",
	"Looks up an gameobject by",
	"==== Команда в игре ====",
	"========================",
	"Поддержка(",
	"Мастер(",
	"Редактор(",
	"Арбитр(",
	"Глава поддержки(",
	"Старший разработчик(",
	"Ведущий(",
	"Рецензент(",
	"Главный редактор(",
	"Рецензент(",
	"Экономист(",
	"Модератор(",
	"Старший рецензент(",
	"Старший ведущий(",
	"Вы получили предмет: ",
	"Группа превращена в рейд.",
	"Ежедневные задания обновились!",
	"Incorrect values.",
	"- Mailbox",
	"│", "├─",
	"Ваша группа расформирована",
	"(self)",
	"Creature moved.",
	"NPC Flags:",
	"Flags Extra:",
	"Armor:",
	"InstanceID:",
	"Loot:",
	"Dynamic Flags:",
	"Unit Flags ",
	"Unit Flags: ",
	"Health ",
	"NPC currently selected by player",
	"SpawnID: ",
	"Faction: ",
	"DisplayID: ",
	"Compatibility Mode: ",
	"Level: ",
	"EquipmentId: ",
	"Target unit has",
	"MechanicImmuneMask: ",
	"UNIT_FLAG",
	"Отчет (",
	".m2", ".wmo",
	"There is no such subcommand",
	"You should select a creature",
	"Result limit reached",
	"Invalid item",
	"Mail sent to",
	"cff7fc80011213",
	"RBAC data reloaded",
	"Incorrect syntax",
	"Modify the hp of the selected player",
}

func createCleanedLog(inputFile multipart.File) (*os.File, error) {
	scanner := bufio.NewScanner(inputFile)
	var removed int
	f, err := ioutil.TempFile("", "output-")
	if err != nil {
		return nil, errors.New("unable to create temp file on server")
	}
mainLoop:
	for scanner.Scan() {
		line := scanner.Text()

		for _, rule := range ignoreRules {
			if strings.Contains(line, rule) {
				removed++
				continue mainLoop
			}
		}
		_, err = f.WriteString(line + "\n")
		if err != nil {
			return nil, errors.New("unable to write cleaned log to temp file")
		}

	}
	_, err = f.WriteString(fmt.Sprintf("\n\nОчищено при помощи dm.rolevik.site. Удалено %v строк.", removed))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f, nil
}

func CleanLog(c *gin.Context) {
	file, err := c.FormFile("input")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !strings.HasSuffix(file.Filename, ".txt") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong file format!"})
		return
	}
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Max file size is 100mb"})
		return
	}
	textFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cleanedFile, err := createCleanedLog(textFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "utf8")
	c.Header("Content-Disposition", "attachment; filename="+cleanedFile.Name())
	c.Header("Content-Type", "application/octet-stream")
	c.FileAttachment(cleanedFile.Name(), "output.txt")
	defer textFile.Close()
	defer os.Remove(cleanedFile.Name())
}
