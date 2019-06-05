use super::*;
// use git2::*;
// use std::path::Path;

pub fn compute_churn(repo: &Repository, commit: &Commit, _interactive: bool) -> Stats {
    let commit_tree = commit.tree().expect("Could not get tree");
    let commit_str: String = commit.id().to_string();

    let mut commit_stats = CommitStats::new();
    commit_stats.merge_commit = commit.parents().count() > 1;
    for p in commit.parents() {
        let parent_tree = p.tree().expect("Failed find parent tree");
        let mut diff_opts = init_diff_opts();
        // let mut blame_opts = init_blame_opts(commit.id());
        let diff = repo
            .diff_tree_to_tree(Some(&parent_tree), Some(&commit_tree), Some(&mut diff_opts))
            .expect("Failed to diff tree to tree");

        let mut line_cb = |d: DiffDelta, _o: Option<DiffHunk>, l: DiffLine| -> bool {
            let path_str = new_path_str_from(&d);
            // let old_path = d
            //     .old_file()
            //     .path()
            //     .expect("Could not find old file in commit.");
            // let blame = repo //SLOW! Redo the blame for every line?!?!
            //     .blame_file(&old_path, Some(&mut blame_opts))
            //     .expect("Failed to execute blame");
            commit_stats
                .commit_path_stats
                .entry(path_str)
                .and_modify(|cps| {
                    // println!(
                    //     "{}\t{}",
                    //     blame
                    //         .get_line(usize::try_from(l.old_lineno().unwrap()).unwrap())
                    //         .unwrap()
                    //         .orig_commit_id(),
                    //     l.origin()
                    // );
                    match l.origin() {
                        '+' => cps.insertions += 1,
                        '-' => cps.deletions += 1,
                        _ => (),
                    }
                })
                .or_insert(CommitPathStats::new());
            true
        };

        diff.foreach(&mut |_, _| -> bool { true }, None, None, Some(&mut line_cb))
            .expect("Failed to iterate over diff");

        // let path = Path::new("src/main.rs");
        // let blame = repo
        //     .blame_file(&path, Some(&mut blame_opts))
        //     .expect("Failed to execute blame");
        // for hunk in blame.iter() {
        // println!(
        //     "Hunk {} Start line: {}",
        //     hunk.orig_commit_id(),
        //     hunk.final_start_line()
        // );
        // }
    }
    let mut stats = Stats::new();
    stats.commit_stats.insert(commit_str, commit_stats);
    return stats;
}

fn init_diff_opts() -> DiffOptions {
    let mut diff_opts = DiffOptions::new();
    diff_opts.ignore_whitespace(true);
    diff_opts.context_lines(0);
    diff_opts.ignore_filemode(true);
    diff_opts.indent_heuristic(true);
    return diff_opts;
}

// fn init_blame_opts(id: Oid) -> BlameOptions {
//     let mut blame_opts = BlameOptions::new();
//     blame_opts.newest_commit(id);
//     blame_opts.track_copies_same_file(true);
//     blame_opts.track_copies_same_commit_moves(false);
//     blame_opts.track_copies_same_commit_copies(false);
//     blame_opts.track_copies_any_commit_copies(false);
//     return blame_opts;
// }

fn new_path_str_from(d: &DiffDelta) -> String {
    let path = d
        .new_file()
        .path()
        .expect("Could not find new file in delta");
    return String::from(path.to_str().expect("Could not convert to path string."));
}

// fn old_path_from(d: &DiffDelta) -> Path {
//     return d
//         .old_file()
//         .path()
//         .expect("Could not find old file in commit.");
// }

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_goofing_stats() {
        let repo = git2::Repository::open(".").unwrap();
        let commit = repo
            .revparse_single("test-goofing")
            .unwrap()
            .peel_to_commit()
            .unwrap();
        let stats = compute_churn(&repo, &commit, false);
        let actual = stats
            .commit_stats
            .get("79caa008ba1f9d06b34b4acc7c03d7fade185a63")
            .unwrap()
            .commit_path_stats
            .get("src/main.rs")
            .unwrap();
        println!("{:#?}", stats);
        assert_eq!(10, actual.insertions);
        assert_eq!(1, actual.deletions);
    }

    #[test]
    fn test_file_origin() {
        let repo = git2::Repository::open(".").unwrap();
        let commit = repo
            .revparse_single("test-file-origin")
            .unwrap()
            .peel_to_commit()
            .unwrap();
        let stats = compute_churn(&repo, &commit, false);
        let actual = stats
            .commit_stats
            .get("6255cfe24e726c0d9222075879e7a2676ac1b5a1")
            .unwrap()
            .commit_path_stats
            .get("testdata/file.txt")
            .unwrap();
        println!("{:#?}", stats);
        assert_eq!(3, actual.insertions);
        assert_eq!(0, actual.deletions);
    }

}
