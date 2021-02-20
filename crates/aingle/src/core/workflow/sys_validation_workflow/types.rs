use super::*;

#[derive(Debug)]
/// The outcome of sys validation
pub(super) enum Outcome {
    /// Moves to app validation
    Accepted,
    /// Moves straight to integration
    SkipAppValidation,
    /// Stays in limbo because another DgdOp
    /// dependency needs to be validated first
    AwaitingOpDep(AnyDgdHash),
    /// Stays in limbo because a dependency could not
    /// be found currently on the DGD.
    /// Note this is not proof it doesn't exist.
    MissingDgdDep,
    /// Moves to integration with status rejected
    Rejected,
}
