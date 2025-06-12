package gognitive

// ContentNode represents a single piece of structured content from a lifelog.
type ContentNode struct {
	Type              string        `json:"type"`                // heading1, heading2, blockquote, etc.
	Content           string        `json:"content"`             // The actual text
	StartTime         string        `json:"startTime,omitempty"` // ISO 8601
	EndTime           string        `json:"endTime,omitempty"`
	StartOffsetMs     int           `json:"startOffsetMs,omitempty"` // ms from start of entry
	EndOffsetMs       int           `json:"endOffsetMs,omitempty"`
	SpeakerName       *string       `json:"speakerName,omitempty"`       // Optional speaker name
	SpeakerIdentifier *string       `json:"speakerIdentifier,omitempty"` // e.g., "user"
	Children          []ContentNode `json:"children,omitempty"`          // Optional nested nodes
}

// Lifelog represents a single lifelog entry.
type Lifelog struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Markdown  string        `json:"markdown,omitempty"` // Full markdown version
	Contents  []ContentNode `json:"contents"`           // Structured blocks
	StartTime string        `json:"startTime"`          // ISO 8601
	EndTime   string        `json:"endTime"`            // ISO 8601
}

// MetaLifelogs holds pagination info.
type MetaLifelogs struct {
	NextCursor string `json:"nextCursor,omitempty"`
	Count      int    `json:"count"`
}

// Meta wraps metadata for lifelog responses.
type Meta struct {
	Lifelogs MetaLifelogs `json:"lifelogs"`
}

// LifelogsResponseData contains the list of lifelogs.
type LifelogsResponseData struct {
	Lifelogs []Lifelog `json:"lifelogs"`
}

// LifelogsResponse is the top-level API response object.
type LifelogsResponse struct {
	Data LifelogsResponseData `json:"data"`
	Meta Meta                 `json:"meta"`
}

// EnrichedLifelog represents a full lifelog entry with tags
type EnrichedLifelog struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	StartTime string        `json:"start_time"`
	EndTime   string        `json:"end_time"`
	Markdown  string        `json:"markdown"`
	Contents  []ContentNode `json:"contents"`
	Tags      []string      `json:"tags"`
}
