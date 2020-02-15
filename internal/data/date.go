package data

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"encoding/json"
)

type (
	// Date stores all data for a given day
	Date struct {
		ID          string       // The identifier for this date
		Word        string       // The word of the day
		Leaderboard Leaderboard  // The leaderboard for the day
		Suggestions []Suggestion // The list of suggested words for the following day
	}
)

const dateFormat = "2006-01-02"

// words were taken from the original inspiration for this app, https://hryanjones.com/guess-my-word/
// That project took the words from 1-1,000 common English words on TV and movies https://en.wiktionary.org/wiki/Wiktionary:Frequency_lists/TV/2006/1-1000
var words []string = []string{"course", "against", "ready", "daughter", "work", "friends", "minute", "though", "supposed", "honey", "point", "start", "check", "alone", "matter", "office", "hospital", "three", "already", "anyway", "important", "tomorrow", "almost", "later", "found", "trouble", "excuse", "hello", "money", "different", "between", "every", "party", "either", "enough", "year", "house", "story", "crazy", "mind", "break", "tonight", "person", "sister", "pretty", "trust", "funny", "gift", "change", "business", "train", "under", "close", "reason", "today", "beautiful", "brother", "since", "bank", "yourself", "without", "until", "forget", "anyone", "promise", "happy", "bake", "worry", "school", "afraid", "cause", "doctor", "exactly", "second", "phone", "look", "feel", "somebody", "stuff", "elephant", "morning", "heard", "world", "chance", "call", "watch", "whatever", "perfect", "dinner", "family", "heart", "least", "answer", "woman", "bring", "probably", "question", "stand", "truth", "problem", "patch", "pass", "famous", "true", "power", "cool", "last", "fish", "remote", "race", "noon", "wipe", "grow", "jumbo", "learn", "itself", "chip", "print", "young", "argue", "clean", "remove", "flip", "flew", "replace", "kangaroo", "side", "walk", "gate", "finger", "target", "judge", "push", "thought", "wear", "desert", "relief", "basic", "bright", "deal", "father", "machine", "know", "step", "exercise", "present", "wing", "lake", "beach", "ship", "wait", "fancy", "eight", "hall", "rise", "river", "round", "girl", "winter", "speed", "long", "oldest", "lock", "kiss", "lava", "garden", "fight", "hook", "desk", "test", "serious", "exit", "branch", "keyboard", "naked", "science", "trade", "quiet", "home", "prison", "blue", "window", "whose", "spot", "hike", "laptop", "dark", "create", "quick", "face", "freeze", "plug", "menu", "terrible", "accept", "door", "touch", "care", "rescue", "ignore", "real", "title", "city", "fast", "season", "town", "picture", "tower", "zero", "engine", "lift", "respect", "time", "mission", "play", "discover", "nail", "half", "unusual", "ball", "tool", "heavy", "night", "farm", "firm", "gone", "help", "easy", "library", "group", "jungle", "taste", "large", "imagine", "normal", "outside", "paper", "nose", "long", "queen", "olive", "doing", "moon", "hour", "protect", "hate", "dead", "double", "nothing", "restaurant", "reach", "note", "tell", "baby", "future", "tall", "drop", "speak", "rule", "pair", "ride", "ticket", "game", "hair", "hurt", "allow", "oven", "live", "horse", "bottle", "rock", "public", "find", "garage", "green", "heat", "plan", "mean", "little", "spend", "nurse", "practice", "wish", "uncle", "core", "stop", "number", "nest", "magazine", "pool", "message", "active", "throw", "pull", "level", "wrist", "bubble", "hold", "movie", "huge", "ketchup", "finish", "pilot", "teeth", "flag", "head", "private", "together", "jewel", "child", "decide", "listen", "garbage", "jealous", "wide", "straight", "fall", "joke", "table", "spread", "laundry", "deep", "quit", "save", "worst", "email", "glass", "scale", "safe", "path", "camera", "excellent", "place", "zone", "luck", "tank", "sign", "report", "myself", "knee", "need", "root", "light", "sure", "page", "life", "space", "magic", "size", "tape", "food", "wire", "period", "mistake", "full", "paid", "horrible", "special", "hidden", "rain", "field", "kick", "ground", "screen", "risky", "junk", "juice", "human", "nobody", "mall", "bathroom", "high", "class", "street", "cold", "metal", "nervous", "bike", "internet", "wind", "lion", "summer", "president", "empty", "square", "jersey", "worm", "popular", "loud", "online", "something", "photo", "knot", "mark", "zebra", "road", "storm", "grab", "record", "said", "floor", "theater", "kitchen", "action", "equal", "nice", "dream", "sound", "fifth", "comfy", "talk", "police", "draw", "bunch", "idea", "jerk", "copy", "success", "team", "favor", "open", "neat", "whale", "gold", "free", "mile", "lying", "meat", "nine", "wonderful", "hero", "quilt", "info", "radio", "move", "early", "remember", "understand", "month", "everyone", "quarter", "center", "universe", "name", "zoom", "inside", "label", "yell", "jacket", "nation", "support", "lunch", "twice", "hint", "jiggle", "boot", "alive", "build", "date", "room", "fire", "music", "leader", "rest", "plant", "connect", "land", "body", "belong", "trick", "wild", "quality", "band", "health", "website", "love", "hand", "okay", "yeah", "dozen", "glove", "give", "thick", "flow", "project", "tight", "join", "cost", "trip", "lower", "magnet", "parent", "grade", "angry", "line", "rich", "owner", "block", "shut", "neck", "write", "hotel", "danger", "impossible", "illegal", "show", "come", "want", "truck", "click", "chocolate", "none", "done", "bone", "hope", "share", "cable", "leaf", "water", "teacher", "dust", "orange", "handle", "unhappy", "guess", "past", "frame", "knob", "winner", "ugly", "lesson", "bear", "gross", "midnight", "grass", "middle", "birthday", "rose", "useless", "hole", "drive", "loop", "color", "sell", "unfair", "send", "crash", "knife", "wrong", "guest", "strong", "weather", "kilometer", "undo", "catch", "neighbor", "stream", "random", "continue", "return", "begin", "kitten", "thin", "pick", "whole", "useful", "rush", "mine", "toilet", "enter", "wedding", "wood", "meet", "stolen", "hungry", "card", "fair", "crowd", "glow", "ocean", "peace", "match", "hill", "welcome", "across", "drag", "island", "edge", "great", "unlock", "feet", "iron", "wall", "laser", "fill", "boat", "weird", "hard", "happen", "tiny", "event", "math", "robot", "recently", "seven", "tree", "rough", "secret", "nature", "short", "mail", "inch", "raise", "warm", "gentle", "glue", "roll", "search", "regular", "here", "count", "hunt", "keep", "week"}

// NewDate will generate a new Date object for the given day
// It is general best practice to first attempt to LoadDate before generating a new one
func NewDate(date time.Time) *Date {
	word, err := generateWord(time.Now())
	if err != nil {
		log.Panicf("Unable to generate new date: %s", err)
	}

	return &Date{
		ID:          date.Format(dateFormat),
		Word:        word,
		Leaderboard: Leaderboard{},
		Suggestions: []Suggestion{},
	}
}

// LoadDate will attempt to load the given date from the data store
func LoadDate(date time.Time) (dte *Date, err error) {
	data, err := db.Get(date.Format(dateFormat))
	if err != nil {
		return nil, fmt.Errorf("Could not load date: %w", err)
	}

	dte = &Date{}
	err = json.Unmarshal(data, dte)
	return
}

// Save will save the current Date to the backend
func (d *Date) Save() error {
	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	if err := enc.Encode(d); err != nil {
		return err
	}

	err := db.Set(d.ID, data.Bytes())
	if err != nil {
		return fmt.Errorf("Unable to save date: %w", err)
	}

	return nil
}
