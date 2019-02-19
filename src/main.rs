use git2::Repository;
use git2::DiffOptions;
use git2::DiffFormat;

fn main() {
    let repo = match Repository::open(".") {
        Ok(repo) => repo,
        Err(e) => panic!("failed to open: {}", e),
    };
    let head = match repo.revparse_single("test-first-test") {
        Ok(head) => head,
        Err(e) => panic!("failed to revparse: {}", e),
    };
    let commit  = head.peel_to_commit().unwrap();
    let commit_tree = match head.peel_to_commit().unwrap().tree() {
        Ok(commit_tree) => commit_tree,
        Err(e) => panic!("failed to revparse: {}", e),
    };
    println!("Commit message is: {}",commit.message().unwrap_or(""));
    for p in commit.parents(){
        let p_tree = match p.tree() {
            Ok(p_tree) => p_tree,
            Err(e) => panic!("failed to revparse: {}", e),
        };
        let mut diff_opts = DiffOptions::new();
        let diff = repo.diff_tree_to_tree(Some(&p_tree), Some(&commit_tree), Some(&mut diff_opts));

        // diff.unwrap().print(DiffFormat::Patch,());
        // println!("Diff: {:#?}", diff);
    }
}

#[cfg(test)]
mod tests {

    #[test]
    fn test_goofing_tag() {
        let repo = match git2::Repository::open(".") {
            Ok(repo) => repo,
            Err(e) => panic!("failed to open: {}", e),
        };
        let head = match repo.revparse_single("test-goofing") {
            Ok(head) => head,
            Err(e) => panic!("failed to revparse: {}", e),
        };

        assert_eq!("79caa008ba1f9d06b34b4acc7c03d7fade185a63", head.id().to_string());
    }
}
