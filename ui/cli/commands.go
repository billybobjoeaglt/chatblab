package cli

import (
	"net"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/billybobjoeaglt/chatlab/common"
	"github.com/billybobjoeaglt/chatlab/config"
)

// XXX: DEPRICATED

type CommandCallback func(line string, args []string)

type Command struct {
	Regex    *regexp.Regexp
	Command  string
	Desc     string
	Args     string
	Example  []string
	Callback CommandCallback
}

var commandArr = []Command{
	Command{
		Regex:   regexp.MustCompile(`\/connect ([^ ]+) ?(.*)`),
		Command: "connect",
		Desc:    "connects to a peer",
		Args:    "/connect [IP] (port)",
		Example: []string{
			"/connect localhost",
			"/connect 192.160.1.24 8080",
		},
		Callback: func(line string, args []string) {
			var port string
			if args[1] != "" {
				port = args[1]
			} else {
				port = strconv.Itoa(common.DefaultPort)
			}
			createConnFunc(net.JoinHostPort(args[0], port))
		},
	},
	Command{
		Regex:   regexp.MustCompile(`\/chats`),
		Command: "chats",
		Desc:    "lists chats",
		Args:    "/chats",
		Example: []string{
			"/chats",
		},
		Callback: func(line string, args []string) {
			logger.Println("--- CHATS ---")
			for name, chat := range chatMap {
				if len(chat) > 0 {
					var printStr string
					printStr += "• " + name
					if len(chat) > 1 {
						printStr += " " + styles["group"]("(group)")
					}
					logger.Println(printStr)
				}
			}
		},
	},
	Command{
		Regex:   regexp.MustCompile(`\/current`),
		Command: "current",
		Desc:    "displays current chat",
		Args:    "/current",
		Example: []string{
			"/current",
		},
		Callback: func(line string, args []string) {
			var printString string
			printString += "Current Chat: " + currentChat + "\n"
			if currentChat != "" {
				printString += "Users: " + strings.Join(chatMap[currentChat], ", ")
			}
			logger.Println(printString)
		},
	},
	Command{
		Regex:   regexp.MustCompile(`\/chat (.+)`),
		Command: "chat",
		Desc:    "switches to the given chat",
		Args:    "/chat [chat name]",
		Example: []string{
			"/chat slaidan_lt",
			"/chat leijurv",
		},
		Callback: func(line string, args []string) {
			chat := args[0]
			hasChat := false
			for name := range chatMap {
				if name == chat {
					hasChat = true
					logger.Println("Connecting to chat", chat)
					currentChat = name
				}
			}
			if !hasChat {
				logger.Println("Error: Missing Chat")
			}
		},
	},
	Command{
		Regex:   regexp.MustCompile(`\/group ([^ ]+) (.+)`),
		Command: "group",
		Desc:    "creates a group",
		Args:    "/group [name] [users,here]",
		Example: []string{
			"/chat slaidan_lt, leijurv",
			"/chat leijurv",
		},
		Callback: func(line string, args []string) {
			groupName := strings.TrimSpace(args[0])
			usersArr := strings.Split(args[1], ",")
			for i := range usersArr {
				usersArr[i] = strings.TrimSpace(usersArr[i])
			}
			for name, arr := range chatMap {
				if name == groupName {
					logger.Println("Error: Chat Already Exists by That Name")
					return
				}
				if reflect.DeepEqual(usersArr, arr) {
					logger.Println("Error: Chat With the Same Users Already Exists")
					return
				}
			}
			AddGroup(groupName, usersArr)
			currentChat = groupName
		},
	},
	Command{
		Regex:   regexp.MustCompile(`\/user (.+)`),
		Command: "user",
		Desc:    "adds a user",
		Args:    "/user [username]",
		Example: []string{
			"/user slaidan_lt",
			"/user leijurv",
		},
		Callback: func(line string, args []string) {
			username := strings.TrimSpace(args[0])
			userArr := []string{username}
			for name, arr := range chatMap {
				if name == username {
					logger.Println("Error: User Already Exists by That Name")
					return
				}
				if reflect.DeepEqual(userArr, arr) {
					logger.Println("Error: Chat With the User Already Exists")
					return
				}
			}
			ok, err := common.DoesUserExist(username)
			if err != nil {
				logger.Println("Error:", err.Error())
				return
			}
			if !ok {
				logger.Println("Error: User Does Not Exist")
				return
			}
			AddUser(username)
		},
	},
	Command{
		Regex:   regexp.MustCompile(`\/settings ([^ ]+) (.+)`),
		Command: "settings",
		Desc:    "changes a setting",
		Args:    "/settings [setting to change] [value]",
		Example: []string{
			"/settings key",
			"/settings password",
		},
		Callback: func(line string, args []string) {
			switch args[0] {
			case "username":
				config.GetConfig().Username = args[1]
			case "key":
				config.GetConfig().PrivateKey = args[1]
			case "save-key":
				if args[1] == "y" {
					pkfp := filepath.Join(common.ProgramDir, "private.key")
					err := common.CopyFile(config.GetConfig().PrivateKey, pkfp)
					if err != nil {
						logger.Println(styles["error"]("Error: " + err.Error()))
						return
					}
					config.GetConfig().PrivateKey = pkfp
				}
			case "password":
				config.Password = args[1]
			case "save-password":
				if args[1] == "y" {
					config.GetConfig().Password = config.Password
				} else if args[1] == "N" {
					config.GetConfig().Password = ""
				} else {
					logger.Println(styles["error"]("Error: Unknown value. Answer y/N"))
				}
			}
		},
	},
}
