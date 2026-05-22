// Package auditor provides audit trail functionality for envdiff comparisons.
//
// It records structured entries each time a diff is performed, capturing:
//   - The timestamp of the operation
//   - The actor (user or system) that triggered it
//   - The files compared
//   - The full diff results
//   - An aggregate summary (total, missing, conflicts, matching)
//
// Entries are persisted as JSON and can be appended across multiple runs,
// forming a chronological log suitable for compliance or debugging.
//
// Example usage:
//
//	entry := auditor.Record(
//		[]string{".env.staging", ".env.prod"},
//		results,
//		auditor.Options{Actor: "deploy-bot"},
//	)
//	if err := auditor.Save("audit.json", entry); err != nil {
//		log.Fatal(err)
//	}
package auditor
