use git2::Repository;

fn main() {
    let repo = match git2::Repository::open(".") {
        Ok(repo) => repo,
        Err(e) => panic!("failed to open: {}", e),
    };
    // println!("Repo head is: {}",head.id());
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
