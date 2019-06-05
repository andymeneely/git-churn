use actix_web::{web, App, HttpServer, HttpResponse, Responder};
use git2::*;

mod git_churn;

fn index(args: web::Path<(String, u32)>) -> impl Responder {
    let repo = Repository::open(".").expect("Could not open repo");
    let commit_str = &args.0;
    let commit = repo
        .revparse_single(&commit_str)
        .expect("Failed revparse")
        .peel_to_commit()
        .expect("Could not peel to commit");
    let stats = git_churn::compute_churn(&repo, &commit, false);
    return HttpResponse::Ok().json(stats);
}

fn main() -> std::io::Result<()> {
    HttpServer::new(|| App::new().service(web::resource("/dogfood/{commit}/{id}").to(index)))
        .bind("127.0.0.1:8080")?
        .run()
}
