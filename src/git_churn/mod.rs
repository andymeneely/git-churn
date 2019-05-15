use git2::*;

pub struct ChurnStats {
    log_str: String,
    commits: usize,
    insertions: usize,
    deletions: usize,
    files_changed: usize,
    authors: usize,
    committers: usize,
    merge_commits: usize,
    merges: usize,
    self_deletions: Option<usize>,
    authors_affected_by_other_deletions: Option<usize>,
}

impl ChurnStats {
    pub fn pretty_print(&self) {
        println!(
            r#"
Git Churn stats for {log_str}
  {commits:<3} commits
  {insertions:<3} insertions
  {deletions:<3} deletions
  {files_changed:<3} files_changed
  {authors:<3} authors
  {committers:<3} committers
  {merge_commits:<3} merge_commits
  {merges:<3} merges
  {self_deletions:?} self deletions
  {authors_affected_by_other_deletions:?} authors affected by other deletions
"#,
            log_str = self.log_str,
            commits = self.commits,
            insertions = self.insertions,
            deletions = self.deletions,
            files_changed = self.files_changed,
            authors = self.authors,
            committers = self.committers,
            merge_commits = self.merge_commits,
            merges = self.merges,
            self_deletions = self.self_deletions,
            authors_affected_by_other_deletions = self.authors_affected_by_other_deletions,
        );
    }

    pub fn new() -> ChurnStats {
        ChurnStats {
            log_str: String::from("TEST log"),
            commits: 0,
            insertions: 0,
            deletions: 0,
            files_changed: 0,
            authors: 0,
            committers: 0,
            merge_commits: 0,
            merges: 0,
            self_deletions: None,
            authors_affected_by_other_deletions: None,
        }
    }

    pub fn set(&mut self, diff_stats: DiffStats) -> &mut ChurnStats {
        self.insertions = diff_stats.insertions();
        self.deletions = diff_stats.deletions();
        self.files_changed = diff_stats.files_changed();
        return self;
    }
}
