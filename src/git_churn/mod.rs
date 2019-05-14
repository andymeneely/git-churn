pub mod git_churn {
    pub struct ChurnStats {
        commits: u32,
        insertions: u32,
        deletions: u32,
        authors: u32,
        committers: u32,
        merge_commits: u32,
        merges: u32,
        self_deletions: Option<u32>,
        authors_affected_by_other_deletions: Option<u32>,
    }

    impl ChurnStats {
        //TODO to_json string
        // fn to_json() -> String {
        //
        // }
    }
}
