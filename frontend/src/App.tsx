import * as React from "react";
import { Route, BrowserRouter as Router, Switch } from "react-router-dom";
import AppContext from "./AppContext";
import Details from "./components/details/Details";
import Home from "./components/home/Home";
import Layout from "./components/layout/Layout";
import Portal from "./components/portal/Portal";
import Getdebt from "./components/getdebt/getdebt";
import Bank from "./components/bank/bank";
import Company from "./components/company/company";
import Tags from "./components/tags/Tags";
import useHttp from "./hooks/http";
import User1 from "./components/user/user";
import UserInfo from "./models/UserInfo";
import Post from "./components/post/post"
import Page from "./components/page/Page"

const App: React.FunctionComponent = () => {
  const [user, setUser] = React.useState<UserInfo>();
  const [selectedKey, setSelectedKey] = React.useState<string>();
  const userInfoRequest = useHttp<UserInfo>("/api/user/self", "GET");

  React.useEffect(() => {
    if (userInfoRequest.data && !userInfoRequest.loading) {
      setUser(userInfoRequest.data);
    }
  }, [userInfoRequest.data, userInfoRequest.loading]);

  React.useEffect(() => {
    userInfoRequest.fire();
  }, []);

  return (
    <AppContext.Provider value={{ user, setUser, selectedKey, setSelectedKey }}>
      <Router>
        <Layout>
          <Switch>
            <Route exact path="/">
              <Home />
            </Route>
            <Route path="/tags">
              <Tags />
            </Route>
            <Route path="/portal">
              <Portal />
            </Route>

            <Route path="/post">
              <Post />
            </Route>
            <Route path="/getdebt">
              <Getdebt />
            </Route>
            <Route path="/company">
              <Company />
            </Route>
            <Route path="/bank">
              <Bank />
            </Route>
            <Route path="/user">
              <User1 />
            </Route>
            <Route path="/details/:id">
              <Details />
            </Route>
            <Route path="/page/:id">
              <Page />
            </Route>
          </Switch>
        </Layout>
      </Router>
    </AppContext.Provider>
  );
};

export default App;