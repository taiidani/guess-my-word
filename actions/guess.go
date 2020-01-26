package actions

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
)

type guess struct {
	word  string
	start time.Time
}

type guessReply struct {
	Guess   string `json:"guess"`
	Correct bool   `json:"correct"`
	After   bool   `json:"after"`
	Before  bool   `json:"before"`
	Error   string `json:"error"`
}

var (
	// words were taken from the original inspiration for this app, https://hryanjones.com/guess-my-word/
	// That project took the words from 1-1,000 common English words on TV and movies https://en.wiktionary.org/wiki/Wiktionary:Frequency_lists/TV/2006/1-1000
	words []string = []string{"course", "against", "ready", "daughter", "work", "friends", "minute", "though", "supposed", "honey", "point", "start", "check", "alone", "matter", "office", "hospital", "three", "already", "anyway", "important", "tomorrow", "almost", "later", "found", "trouble", "excuse", "hello", "money", "different", "between", "every", "party", "either", "enough", "year", "house", "story", "crazy", "mind", "break", "tonight", "person", "sister", "pretty", "trust", "funny", "gift", "change", "business", "train", "under", "close", "reason", "today", "beautiful", "brother", "since", "bank", "yourself", "without", "until", "forget", "anyone", "promise", "happy", "bake", "worry", "school", "afraid", "cause", "doctor", "exactly", "second", "phone", "look", "feel", "somebody", "stuff", "elephant", "morning", "heard", "world", "chance", "call", "watch", "whatever", "perfect", "dinner", "family", "heart", "least", "answer", "woman", "bring", "probably", "question", "stand", "truth", "problem", "patch", "pass", "famous", "true", "power", "cool", "last", "fish", "remote", "race", "noon", "wipe", "grow", "jumbo", "learn", "itself", "chip", "print", "young", "argue", "clean", "remove", "flip", "flew", "replace", "kangaroo", "side", "walk", "gate", "finger", "target", "judge", "push", "thought", "wear", "desert", "relief", "basic", "bright", "deal", "father", "machine", "know", "step", "exercise", "present", "wing", "lake", "beach", "ship", "wait", "fancy", "eight", "hall", "rise", "river", "round", "girl", "winter", "speed", "long", "oldest", "lock", "kiss", "lava", "garden", "fight", "hook", "desk", "test", "serious", "exit", "branch", "keyboard", "naked", "science", "trade", "quiet", "home", "prison", "blue", "window", "whose", "spot", "hike", "laptop", "dark", "create", "quick", "face", "freeze", "plug", "menu", "terrible", "accept", "door", "touch", "care", "rescue", "ignore", "real", "title", "city", "fast", "season", "town", "picture", "tower", "zero", "engine", "lift", "respect", "time", "mission", "play", "discover", "nail", "half", "unusual", "ball", "tool", "heavy", "night", "farm", "firm", "gone", "help", "easy", "library", "group", "jungle", "taste", "large", "imagine", "normal", "outside", "paper", "nose", "long", "queen", "olive", "doing", "moon", "hour", "protect", "hate", "dead", "double", "nothing", "restaurant", "reach", "note", "tell", "baby", "future", "tall", "drop", "speak", "rule", "pair", "ride", "ticket", "game", "hair", "hurt", "allow", "oven", "live", "horse", "bottle", "rock", "public", "find", "garage", "green", "heat", "plan", "mean", "little", "spend", "nurse", "practice", "wish", "uncle", "core", "stop", "number", "nest", "magazine", "pool", "message", "active", "throw", "pull", "level", "wrist", "bubble", "hold", "movie", "huge", "ketchup", "finish", "pilot", "teeth", "flag", "head", "private", "together", "jewel", "child", "decide", "listen", "garbage", "jealous", "wide", "straight", "fall", "joke", "table", "spread", "laundry", "deep", "quit", "save", "worst", "email", "glass", "scale", "safe", "path", "camera", "excellent", "place", "zone", "luck", "tank", "sign", "report", "myself", "knee", "need", "root", "light", "sure", "page", "life", "space", "magic", "size", "tape", "food", "wire", "period", "mistake", "full", "paid", "horrible", "special", "hidden", "rain", "field", "kick", "ground", "screen", "risky", "junk", "juice", "human", "nobody", "mall", "bathroom", "high", "class", "street", "cold", "metal", "nervous", "bike", "internet", "wind", "lion", "summer", "president", "empty", "square", "jersey", "worm", "popular", "loud", "online", "something", "photo", "knot", "mark", "zebra", "road", "storm", "grab", "record", "said", "floor", "theater", "kitchen", "action", "equal", "nice", "dream", "sound", "fifth", "comfy", "talk", "police", "draw", "bunch", "idea", "jerk", "copy", "success", "team", "favor", "open", "neat", "whale", "gold", "free", "mile", "lying", "meat", "nine", "wonderful", "hero", "quilt", "info", "radio", "move", "early", "remember", "understand", "month", "everyone", "quarter", "center", "universe", "name", "zoom", "inside", "label", "yell", "jacket", "nation", "support", "lunch", "twice", "hint", "jiggle", "boot", "alive", "build", "date", "room", "fire", "music", "leader", "rest", "plant", "connect", "land", "body", "belong", "trick", "wild", "quality", "band", "health", "website", "love", "hand", "okay", "yeah", "dozen", "glove", "give", "thick", "flow", "project", "tight", "join", "cost", "trip", "lower", "magnet", "parent", "grade", "angry", "line", "rich", "owner", "block", "shut", "neck", "write", "hotel", "danger", "impossible", "illegal", "show", "come", "want", "truck", "click", "chocolate", "none", "done", "bone", "hope", "share", "cable", "leaf", "water", "teacher", "dust", "orange", "handle", "unhappy", "guess", "past", "frame", "knob", "winner", "ugly", "lesson", "bear", "gross", "midnight", "grass", "middle", "birthday", "rose", "useless", "hole", "drive", "loop", "color", "sell", "unfair", "send", "crash", "knife", "wrong", "guest", "strong", "weather", "kilometer", "undo", "catch", "neighbor", "stream", "random", "continue", "return", "begin", "kitten", "thin", "pick", "whole", "useful", "rush", "mine", "toilet", "enter", "wedding", "wood", "meet", "stolen", "hungry", "card", "fair", "crowd", "glow", "ocean", "peace", "match", "hill", "welcome", "across", "drag", "island", "edge", "great", "unlock", "feet", "iron", "wall", "laser", "fill", "boat", "weird", "hard", "happen", "tiny", "event", "math", "robot", "recently", "seven", "tree", "rough", "secret", "nature", "short", "mail", "inch", "raise", "warm", "gentle", "glue", "roll", "search", "regular", "here", "count", "hunt", "keep", "week"}
)

// GuessHandler is an API handler to process a user's guess.
func GuessHandler(c buffalo.Context) error {
	guess := extractGuess(c)
	reply := guessReply{}
	reply.Guess = guess.word

	// Validate the guess
	if len(reply.Guess) == 0 {
		reply.Error = "Guess must not be empty"
	} else if !validateWord(reply.Guess) {
		reply.Error = "Guess must be a valid Scrabble word"
	}

	// Generate the word for the day
	word, err := generateWord(guess.start)
	if err != nil {
		return err
	}

	if reply.Error == "" {
		switch strings.Compare(reply.Guess, word) {
		case -1:
			reply.After = true
		case 1:
			reply.Before = true
		case 0:
			reply.Correct = true
		}
	}

	return c.Render(200, r.JSON(reply))
}

func extractGuess(c buffalo.Context) guess {
	ret := guess{}
	ret.word = strings.ToLower(strings.TrimSpace(c.Param("word")))

	startStr := strings.TrimSpace(c.Param("start"))
	if startUnix, err := strconv.ParseInt(startStr, 10, 64); err == nil {
		ret.start = time.Unix(startUnix, 0)
	}

	return ret
}

func generateWord(seed time.Time) (string, error) {
	if seed.Unix() == 0 {
		return "", fmt.Errorf("Invalid timestamp for word")
	}

	day := seed.UTC()
	return words[(day.Year()*day.YearDay())%len(words)], nil
}

func validateWord(word string) bool {
	for _, line := range scrabble {
		if line == word {
			return true
		}
	}

	return false
}
