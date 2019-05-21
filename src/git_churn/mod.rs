use git2::*;
use std::collections::HashMap;

mod metrics;

pub struct Stats {
    log_str: String,
    commit_stats: HashMap<String, CommitStats>,
    files_changed: usize,
    authors: usize,
    committers: usize,
}

pub struct CommitStats {
    merge_commit: bool,
    commit_path_stats: HashMap<String, CommitPathStats>,
    authors_affected_by_other_deletions: Option<usize>,
}

pub struct CommitPathStats {
    insertions: usize,
    deletions: usize,
    self_deletions: Option<usize>,
    authors_affected_by_other_deletions: Option<usize>,
}

impl Stats {
    pub fn pretty_print(&self) {

//         println!(
//             r#"
// Git Churn stats for {log_str}
//   {commits:<3} commits
//   {insertions:<3} insertions
//   {deletions:<3} deletions
//   {files_changed:<3} files_changed
//   {authors:<3} authors
//   {committers:<3} committers
//   {merge_commits:<3} merge_commits
//   {merges:<3} merges
//   {self_deletions:?} self deletions
//   {authors_affected_by_other_deletions:?} authors affected by other deletions
// "#,
//             log_str = self.log_str,
//             commits = self.commits,
//             insertions = self.insertions,
//             deletions = self.deletions,
//             files_changed = self.files_changed,
//             authors = self.authors,
//             committers = self.committers,
//             merge_commits = self.merge_commits,
//             merges = self.merges,
//             self_deletions = self.self_deletions,
//             authors_affected_by_other_deletions = self.authors_affected_by_other_deletions,
//         );
    }

    pub fn new() -> Stats {
        Stats {
            log_str: String::from(""),
            commit_stats: HashMap::new(),
            files_changed: 0,
            authors: 0,
            committers: 0,
        }
    }
}

pub fn compute_churn(repo: &Repository, commit: &Commit, _interactive: bool) -> Stats {
    return self::metrics::compute_churn(&repo, &commit, false);
}
