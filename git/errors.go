package git

import (
	"regexp"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ErrorCode int64

type GitError struct {
	Stderr    string
	ErrorCode ErrorCode
	Message   string
}

func (g *Git) ParseGitError(stderr string) error {
	e := GitError{
		Stderr: stderr,
	}
	// Determine error code from stderr
	e.parse()
	// Get message based on error code
	msg := getGitErrorMessage(e.ErrorCode)
	if msg != "" {
		e.Message = msg
	} else {
		// If not a standard error, the message will simply be stderr
		e.Message = e.Stderr
	}

	// Handle events for some errors.
	switch e.ErrorCode {
	case RebaseConflicts:
		runtime.EventsEmit(g.AppCtx, "rebase-conflicts")
	case MergeConflicts:
		runtime.EventsEmit(g.AppCtx, "merge-conflicts")
	case RevertConflicts:
		runtime.EventsEmit(g.AppCtx, "revert-conflicts")
	case UnresolvedConflicts:
		runtime.EventsEmit(g.AppCtx, "unresolved-conflicts")
	}

	return &e
}

func (e *GitError) Error() string {
	return e.Message
}

const (
	Undefined ErrorCode = iota
	SSHKeyAuditUnverified
	SSHAuthenticationFailed
	SSHPermissionDenied
	HTTPSAuthenticationFailed
	RemoteDisconnection
	HostDown
	RebaseConflicts
	MergeConflicts
	HTTPSRepositoryNotFound
	SSHRepositoryNotFound
	PushNotFastForward
	BranchDeletionFailed
	DefaultBranchDeletionFailed
	RevertConflicts
	EmptyRebasePatch
	NoMatchingRemoteBranch
	NoExistingRemoteBranch
	NothingToCommit
	NoSubmoduleMapping
	SubmoduleRepositoryDoesNotExist
	InvalidSubmoduleSHA
	LocalPermissionDenied
	InvalidMerge
	InvalidRebase
	NonFastForwardMergeIntoEmptyHead
	PatchDoesNotApply
	BranchAlreadyExists
	BadRevision
	NotAGitRepository
	CannotMergeUnrelatedHistories
	LFSAttributeDoesNotMatch
	BranchRenameFailed
	PathDoesNotExist
	InvalidObjectName
	OutsideRepository
	LockFileAlreadyExists
	NoMergeToAbort
	LocalChangesOverwritten
	UnresolvedConflicts
	GPGFailedToSignData
	ConflictModifyDeletedInBranch
	ConfigLockFileAlreadyExists
	RemoteAlreadyExists
	TagAlreadyExists
	MergeWithLocalChanges
	RebaseWithLocalChanges
	MergeCommitNoMainlineOption
	UnsafeDirectory
	PathExistsButNotInRef
	// Start of GitHub-specific error codes
	PushWithFileSizeExceedingLimit
	HexBranchNameRejected
	ForcePushRejected
	InvalidRefLength
	ProtectedBranchRequiresReview
	ProtectedBranchForcePush
	ProtectedBranchDeleteRejected
	ProtectedBranchRequiredStatus
	PushWithPrivateEmail
	// End of GitHub-specific error codes
)

func (e *GitError) parse() {
	stderr := e.Stderr
	for _, r := range getGitErrorRegexes() {
		if match, _ := regexp.MatchString(r.Regex, stderr); match {
			e.ErrorCode = r.Code
		}
	}
}

type GitErrorRegex struct {
	Code  ErrorCode
	Regex string
}

func getGitErrorRegexes() []GitErrorRegex {
	return []GitErrorRegex{
		{
			Code:  SSHKeyAuditUnverified,
			Regex: "ERROR: ([\\s\\S]+?)\\n+\\[EPOLICYKEYAGE\\]\\n+fatal: Could not read from remote repository.",
		},
		{
			Code:  HTTPSAuthenticationFailed,
			Regex: "fatal: Authentication failed for 'https://",
		},
		{
			Code:  SSHAuthenticationFailed,
			Regex: "fatal: Authentication failed",
		},
		{
			Code:  SSHPermissionDenied,
			Regex: "fatal: Could not read from remote repository.",
		},
		{
			Code:  HTTPSAuthenticationFailed,
			Regex: "The requested URL returned error: 403",
		},
		{
			Code:  RemoteDisconnection,
			Regex: "fatal: [Tt]he remote end hung up unexpectedly",
		},
		{
			Code:  HostDown,
			Regex: "fatal: unable to access '(.+)': Failed to connect to (.+): Host is down",
		},
		{
			Code:  HostDown,
			Regex: "Cloning into '(.+)'...\nfatal: unable to access '(.+)': Could not resolve host: (.+)",
		},
		{
			Code:  RebaseConflicts,
			Regex: "Resolve all conflicts manually, mark them as resolved with",
		},
		{
			Code:  MergeConflicts,
			Regex: "(Merge conflict|Automatic merge failed; fix conflicts and then commit the result)",
		},
		{
			Code:  HTTPSRepositoryNotFound,
			Regex: "fatal: repository '(.+)' not found",
		},
		{
			Code:  SSHRepositoryNotFound,
			Regex: "ERROR: Repository not found",
		},
		{
			Code:  PushNotFastForward,
			Regex: "\\((non-fast-forward|fetch first)\\)\nerror: failed to push some refs to '.*'",
		},
		{
			Code:  BranchDeletionFailed,
			Regex: "error: unable to delete '(.+)': remote ref does not exist",
		},
		{
			Code:  DefaultBranchDeletionFailed,
			Regex: "\\[remote rejected\\] (.+) \\(deletion of the current branch prohibited\\)",
		},
		{
			Code:  RevertConflicts,
			Regex: "error: could not revert .*\nhint: after resolving the conflicts, mark the corrected paths\nhint: with 'git add <paths>' or 'git rm <paths>'\nhint: and commit the result with 'git commit'",
		},
		{
			Code:  EmptyRebasePatch,
			Regex: "Applying: .*\nNo changes - did you forget to use 'git add'\\?\nIf there is nothing left to stage, chances are that something else\n.*",
		},
		{
			Code:  NoMatchingRemoteBranch,
			Regex: "There are no candidates for (rebasing|merging) among the refs that you just fetched.\nGenerally this means that you provided a wildcard refspec which had no\nmatches on the remote end.",
		},
		{
			Code:  NoExistingRemoteBranch,
			Regex: "Your configuration specifies to merge with the ref '(.+)'\nfrom the remote, but no such ref was fetched.",
		},
		{
			Code:  NothingToCommit,
			Regex: "nothing to commit",
		},
		{
			Code:  NoSubmoduleMapping,
			Regex: "[Nn]o submodule mapping found in .gitmodules for path '(.+)'",
		},
		{
			Code:  SubmoduleRepositoryDoesNotExist,
			Regex: "fatal: repository '(.+)' does not exist\nfatal: clone of '.+' into submodule path '(.+)' failed",
		},
		{
			Code:  InvalidSubmoduleSHA,
			Regex: "Fetched in submodule path '(.+)', but it did not contain (.+). Direct fetching of that commit failed.",
		},
		{
			Code:  LocalPermissionDenied,
			Regex: "fatal: could not create work tree dir '(.+)'.*: Permission denied",
		},
		{
			Code:  InvalidMerge,
			Regex: "merge: (.+) - not something we can merge",
		},
		{
			Code:  InvalidRebase,
			Regex: "invalid upstream (.+)",
		},
		{
			Code:  NonFastForwardMergeIntoEmptyHead,
			Regex: "fatal: Non-fast-forward commit does not make sense into an empty head",
		},
		{
			Code:  PatchDoesNotApply,
			Regex: "error: (.+): (patch does not apply|already exists in working directory)",
		},
		{
			Code:  BranchAlreadyExists,
			Regex: "fatal: [Aa] branch named '(.+)' already exists.?",
		},
		{
			Code:  BadRevision,
			Regex: "fatal: bad revision '(.*)'",
		},
		{
			Code:  NotAGitRepository,
			Regex: "fatal: [Nn]ot a git repository \\(or any of the parent directories\\): (.*)",
		},
		{
			Code:  CannotMergeUnrelatedHistories,
			Regex: "fatal: refusing to merge unrelated histories",
		},
		{
			Code:  LFSAttributeDoesNotMatch,
			Regex: "The .+ attribute should be .+ but is .+",
		},
		{
			Code:  BranchRenameFailed,
			Regex: "fatal: Branch rename failed",
		},
		{
			Code:  PathDoesNotExist,
			Regex: "fatal: path '(.+)' does not exist .+",
		},
		{
			Code:  InvalidObjectName,
			Regex: "fatal: invalid object name '(.+)'.",
		},
		{
			Code:  OutsideRepository,
			Regex: "fatal: .+: '(.+)' is outside repository",
		},
		{
			Code:  LockFileAlreadyExists,
			Regex: "Another git process seems to be running in this repository, e.g.",
		},
		{
			Code:  NoMergeToAbort,
			Regex: "fatal: There is no merge to abort",
		},
		{
			Code:  LocalChangesOverwritten,
			Regex: "error: (?:Your local changes to the following|The following untracked working tree) files would be overwritten by checkout:",
		},
		{
			Code:  UnresolvedConflicts,
			Regex: "You must edit all merge conflicts and then\nmark them as resolved using git add|fatal: Exiting because of an unresolved conflict",
		},
		{
			Code:  GPGFailedToSignData,
			Regex: "error: gpg failed to sign the data",
		},
		{
			Code:  ConflictModifyDeletedInBranch,
			Regex: "CONFLICT \\(modify/delete\\): (.+) deleted in (.+) and modified in (.+)",
		},
		{
			Code:  ConfigLockFileAlreadyExists,
			Regex: "error: could not lock config file (.+): File exists",
		},
		{
			Code:  RemoteAlreadyExists,
			Regex: "error: remote (.+) already exists.",
		},
		{
			Code:  TagAlreadyExists,
			Regex: "fatal: tag '(.+)' already exists",
		},
		{
			Code:  MergeWithLocalChanges,
			Regex: "error: Your local changes to the following files would be overwritten by merge:\n",
		},
		{
			Code:  RebaseWithLocalChanges,
			Regex: "error: cannot (pull with rebase|rebase): You have unstaged changes\\.\n\\s*error: [Pp]lease commit or stash them\\.",
		},
		{
			Code:  MergeCommitNoMainlineOption,
			Regex: "error: commit (.+) is a merge but no -m option was given",
		},
		{
			Code:  MergeCommitNoMainlineOption,
			Regex: "fatal: detected dubious ownership in repository at (.+)",
		},
		{
			Code:  UnsafeDirectory,
			Regex: "fatal: detected dubious ownership in repository at (.+)",
		},
		{
			Code:  PathExistsButNotInRef,
			Regex: "fatal: path '(.+)' exists on disk, but not in '(.+)'",
		},
		// GitHub-specific errors
		{
			Code:  PushWithFileSizeExceedingLimit,
			Regex: "error: GH001: ",
		},
		{
			Code:  HexBranchNameRejected,
			Regex: "error: GH002: ",
		},
		{
			Code:  ForcePushRejected,
			Regex: "error: GH003: Sorry, force-pushing to (.+) is not allowed.",
		},
		{
			Code:  InvalidRefLength,
			Regex: "error: GH005: Sorry, refs longer than (.+) bytes are not allowed",
		},
		{
			Code:  ProtectedBranchRequiresReview,
			Regex: "error: GH006: Protected branch update failed for (.+)\nremote: error: At least one approved review is required",
		},
		{
			Code:  ProtectedBranchForcePush,
			Regex: "error: GH006: Protected branch update failed for (.+)\nremote: error: Cannot force-push to a protected branch",
		},
		{
			Code:  ProtectedBranchDeleteRejected,
			Regex: "error: GH006: Protected branch update failed for (.+)\nremote: error: Cannot delete a protected branch",
		},
		{
			Code:  ProtectedBranchRequiredStatus,
			Regex: "error: GH006: Protected branch update failed for (.+).\nremote: error: Required status check \"(.+)\" is expected",
		},
		{
			Code:  PushWithPrivateEmail,
			Regex: "error: GH007: Your push would publish a private email address.",
		},
		// End GitHub-specific errors
	}
}

func getGitErrorMessage(code ErrorCode) string {
	switch code {
	case SSHKeyAuditUnverified:
		return "The SSH key is unverified."
	case RemoteDisconnection:
		return "The remote disconnected. Check your Internet connection and try again."
	case HostDown:
		return "The host is down. Check your Internet connection and try again."
	case RebaseConflicts:
		return "We found some conflicts while trying to rebase. Please resolve the conflicts before continuing."
	case MergeConflicts:
		return "We found some conflicts while trying to merge. Please resolve the conflicts and commit the changes."
	case HTTPSRepositoryNotFound:
	case SSHRepositoryNotFound:
		return "The repository does not seem to exist anymore. You may not have access, or it may have been deleted or renamed."
	case PushNotFastForward:
		return "The repository has been updated since you last pulled. Try pulling before pushing."
	case BranchDeletionFailed:
		return "Could not delete the branch. It was probably already deleted."
	case DefaultBranchDeletionFailed:
		return "The branch is the repository's default branch and cannot be deleted."
	case RevertConflicts:
		return "To finish reverting, please merge and commit the changes."
	case EmptyRebasePatch:
		return "There aren't any changes left to apply."
	case NoMatchingRemoteBranch:
		return "There aren't any remote branches that match the current branch."
	case NothingToCommit:
		return "There are no changes to commit."
	case NoSubmoduleMapping:
		return "A submodule was removed from .gitmodules, but the folder still exists in the repository. Delete the folder, commit the change, then try again."
	case SubmoduleRepositoryDoesNotExist:
		return "A submodule points to a location which does not exist."
	case InvalidSubmoduleSHA:
		return "A submodule points to a commit which does not exist."
	case LocalPermissionDenied:
		return "Permission denied."
	case InvalidMerge:
		return "This is not something we can merge."
	case InvalidRebase:
		return "This is not something we can rebase."
	case NonFastForwardMergeIntoEmptyHead:
		return "The merge you attempted is not a fast-forward, so it cannot be performed on an empty branch."
	case PatchDoesNotApply:
		return "The requested changes conflict with one or more files in the repository."
	case BranchAlreadyExists:
		return "A branch with that name already exists."
	case BadRevision:
		return "Bad revision."
	case NotAGitRepository:
		return "This is not a git repository."
	case ProtectedBranchForcePush:
		return "This branch is protected from force-push operations."
	case ProtectedBranchRequiresReview:
		return "This branch is protected and any changes requires an approved review. Open a pull request with changes targeting this branch instead."
	case PushWithFileSizeExceedingLimit:
		return "The push operation includes a file which exceeds GitHub's file size restriction of 100MB. Please remove the file from history and try again."
	case HexBranchNameRejected:
		return "The branch name cannot be a 40-character string of hexadecimal characters, as this is the format that Git uses for representing objects."
	case ForcePushRejected:
		return "The force push has been rejected for the current branch."
	case InvalidRefLength:
		return "A ref cannot be longer than 255 characters."
	case CannotMergeUnrelatedHistories:
		return "Unable to merge unrelated histories in this repository."
	case PushWithPrivateEmail:
		return "Cannot push these commits as they contain an email address marked as private on GitHub. To push anyway, visit https://github.com/settings/emails, uncheck \"Keep my email address private\", then switch back to GitHub Desktop to push your commits. You can then enable the setting again."
	case LFSAttributeDoesNotMatch:
		return "Git LFS attribute found in global Git configuration does not match expected value."
	case ProtectedBranchDeleteRejected:
		return "This branch cannot be deleted from the remote repository because it is marked as protected."
	case ProtectedBranchRequiredStatus:
		return "The push was rejected by the remote server because a required status check has not been satisfied."
	case BranchRenameFailed:
		return "The branch could not be renamed."
	case PathDoesNotExist:
		return "The path does not exist on disk."
	case InvalidObjectName:
		return "The object was not found in the Git repository."
	case OutsideRepository:
		return "This path is not a valid path inside the repository."
	case LockFileAlreadyExists:
		return "A lock file already exists in the repository, which blocks this operation from completing."
	case NoMergeToAbort:
		return "There is no merge in progress, so there is nothing to abort."
	case NoExistingRemoteBranch:
		return "The remote branch does not exist."
	case LocalChangesOverwritten:
		return "Unable to switch branches as there are working directory changes which would be overwritten. Please commit or stash your changes."
	case UnresolvedConflicts:
		return "There are unresolved conflicts in the working directory."
	case ConfigLockFileAlreadyExists:
		return ""
	case RemoteAlreadyExists:
		return ""
	case TagAlreadyExists:
		return "A tag with that name already exists"
	case MergeWithLocalChanges:
	case RebaseWithLocalChanges:
	case GPGFailedToSignData:
	case ConflictModifyDeletedInBranch:
	case MergeCommitNoMainlineOption:
	case UnsafeDirectory:
	case PathExistsButNotInRef:
		return ""
	}
	return ""
}
