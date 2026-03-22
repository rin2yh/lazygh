package gh

const (
	queryGetReviewContext = `
		query($owner: String!, $name: String!, $number: Int!) {
			repository(owner: $owner, name: $name) {
				pullRequest(number: $number) {
					id
					headRefOid
				}
			}
		}`

	mutationStartPendingReview = `
		mutation($pullRequestId: ID!, $commitOID: GitObjectID!) {
			addPullRequestReview(input: {
				pullRequestId: $pullRequestId,
				commitOID: $commitOID
			}) {
				pullRequestReview { id }
			}
		}`

	mutationAddReviewComment = `
		mutation(
			$pullRequestReviewId: ID!,
			$body: String!,
			$path: String!,
			$line: Int!,
			$side: DiffSide!,
			$startLine: Int,
			$startSide: DiffSide
		) {
			addPullRequestReviewThread(input: {
				pullRequestReviewId: $pullRequestReviewId,
				body: $body,
				path: $path,
				line: $line,
				side: $side,
				startLine: $startLine,
				startSide: $startSide
			}) {
				thread {
					comments(first: 1) {
						nodes { id }
					}
				}
			}
		}`

	mutationDeleteReviewComment = `
		mutation($id: ID!) {
			deletePullRequestReviewComment(input: { id: $id }) {
				clientMutationId
			}
		}`

	mutationUpdateReviewComment = `
		mutation($id: ID!, $body: String!) {
			updatePullRequestReviewComment(input: {
				pullRequestReviewCommentId: $id,
				body: $body
			}) {
				pullRequestReviewComment { id }
			}
		}`

	mutationSubmitReview = `
		mutation($pullRequestReviewId: ID!, $event: PullRequestReviewEvent!, $body: String!) {
			submitPullRequestReview(input: {
				pullRequestReviewId: $pullRequestReviewId,
				event: $event,
				body: $body
			}) {
				pullRequestReview { id }
			}
		}`

	mutationDeleteReview = `
		mutation($pullRequestReviewId: ID!) {
			deletePullRequestReview(input: {
				pullRequestReviewId: $pullRequestReviewId
			}) {
				clientMutationId
			}
		}`
)
