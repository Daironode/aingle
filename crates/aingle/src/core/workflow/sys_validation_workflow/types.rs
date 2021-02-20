use super::*;

#[derive(Debug)]
/// The outcome of sys validation
pub(super) enum Outcome {
    /// Moves to app validation
    Accepted,
    /// Moves straight to integration
    SkipAppValidation,
<<<<<<< HEAD
    /// Stays in limbo because another DgdOp
    /// dependency needs to be validated first
    AwaitingOpDep(AnyDgdHash),
    /// Stays in limbo because a dependency could not
    /// be found currently on the DGD.
    /// Note this is not proof it doesn't exist.
    MissingDgdDep,
=======
    /// Stays in limbo because another DhtOp
    /// dependency needs to be validated first
    AwaitingOpDep(AnyDhtHash),
    /// Stays in limbo because a dependency could not
    /// be found currently on the DHT.
    /// Note this is not proof it doesn't exist.
    MissingDhtDep,
>>>>>>> master
    /// Moves to integration with status rejected
    Rejected,
}
