#include<ainglelib/aingle.hpp>

using namespace aingle;
using std::string;

class hello: public contract {
	public:
	using contract::contract;
	int128_t add(int128_t a, int128_t b) {
		return a + b;
	}

	string testcase(void) {
		return string(R"(
		[
    	    [{"env":{"witness":[]}, "method":"add", "param":"int:1, int:2", "expected":"int:3"}
    	    ]
		]
		)");
	}
};

AINGLE_DISPATCH( hello, (testcase)(add))
