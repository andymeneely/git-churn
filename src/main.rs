use git2::*;

mod git_churn;

fn main() {
    let repo = Repository::open(".").expect("Could not open repo");
    let commit = repo
        .revparse_single("test-first-test")
        .expect("Failed revparse")
        .peel_to_commit()
        .expect("Could not peel to commit");
    git_churn::compute_churn(&repo, &commit, false).pretty_print();
}

#[cfg(test)]
mod tests {

    #[test]
    fn test_goofing_tag() {
        let repo = git2::Repository::open(".").expect("Failed to open git repo .");
        let head = repo
            .revparse_single("test-goofing")
            .expect("Could not parse");
        assert_eq!(
            "79caa008ba1f9d06b34b4acc7c03d7fade185a63",
            head.id().to_string()
        );
    }

}
