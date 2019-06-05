use git2::*;
use serde_derive::Serialize;
use serde_json::json;
use std::collections::HashMap;

mod metrics;

#[derive(Serialize, Debug)]
pub struct Stats {
    log_str: String,
    commit_stats: HashMap<String, CommitStats>,
    files_changed: usize,
    authors: usize,
    committers: usize,
}

#[derive(Serialize, Debug)]
pub struct CommitStats {
    merge_commit: bool,
    commit_path_stats: HashMap<String, CommitPathStats>,
    authors_affected_by_other_deletions: Option<usize>,
}

#[derive(Serialize, PartialEq, Debug)]
pub struct CommitPathStats {
    insertions: usize,
    deletions: usize,
    self_deletions: Option<usize>,
    authors_affected_by_other_deletions: Option<usize>,
}

impl Stats {
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

impl CommitStats {
    pub fn new() -> CommitStats {
        CommitStats {
            merge_commit: false,
            commit_path_stats: HashMap::new(),
            authors_affected_by_other_deletions: None,
        }
    }
}

impl CommitPathStats {
    pub fn new() -> CommitPathStats {
        CommitPathStats {
            insertions: 0,
            deletions: 0,
            self_deletions: None,
            authors_affected_by_other_deletions: None,
        }
    }
}

pub fn compute_churn(repo: &Repository, commit: &Commit, _interactive: bool) -> Stats {
    return self::metrics::compute_churn(&repo, &commit, false);
}
