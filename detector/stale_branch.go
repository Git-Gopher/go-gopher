package detector

// In Github it's stale after 3 months of inactivity.

/**
query {
	repository(owner: "subquery", name: "subql") {
    refs(refPrefix: "refs/heads/", first: 5) {
      nodes {
        target {
            ... on Commit {
              history(first: 1) {
                nodes {
                  messageHeadline
                  committedDate
                }
              }
            }
        }
      }
    }
  }
}
**/
