#include <crow.h>
#include <iostream>
#include <pqxx/pqxx>
using namespace std;
using namespace pqxx;

int main() {
    cout << "hello";

    crow::SimpleApp app;

    CROW_ROUTE(app, "/getData")
    ([]() {
        cout << "hello";

        try {
            connection C("dbname = finni user = david hostaddr = 127.0.0.1 port = 5432");
            if (C.is_open()) {
                cout << "Opened database successfully: " << C.dbname() << endl;
            } else {
                cout << "Can't open database" << endl;
                return "";
            }

            /* Create SQL statement */
            string sql = "SELECT * FROM finni.patients";

            /* Create a non-transactional object. */
            nontransaction N(C);

            /* Execute SQL query */
            result R(N.exec(sql));

            /* List down all the records */
			string result;
            for (result::const_iterator c = R.begin(); c != R.end(); ++c) {
				//cout << "Salary = " << c[4].as<float>() << endl;
				result.append(c[0].as<string>());
            }
            cout << "Operation done successfully" << endl;
            C.close();
        } catch (const std::exception &e) {
            cerr << e.what() << std::endl;
            return "";
        }

        return "";
    });

    CROW_ROUTE(app, "/")
    ([]() {
        cout << "hello";
        return "ping";
    });
    app.port(18080).run();

    return 1;
}
