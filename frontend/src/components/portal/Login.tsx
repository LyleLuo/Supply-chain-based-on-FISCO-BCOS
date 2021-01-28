import * as React from "react";
import { PrimaryButton, DefaultButton, Stack, Text, TextField } from "@fluentui/react";
import AppContext from "../../AppContext";
import Register from "./Register";
import useHttp from "../../hooks/http";
import UserInfo from "../../models/UserInfo";

const Login: React.FunctionComponent = () => {
  const { setUser } = React.useContext(AppContext);
  const [name, setName] = React.useState<string>();
  const [password, setPassword] = React.useState<string>();
  const [type, setType] = React.useState("login");
  const [error, setError] = React.useState<string>();
  const loginRequest = useHttp<{ status: string; }>("/api/user/login", "POST");
  const userInfoRequest = useHttp<UserInfo>("/api/user/self", "GET");

  React.useEffect(() => {
    if (userInfoRequest.data && !userInfoRequest.loading && setUser) {
      console.log(userInfoRequest.data);
      setUser({
        address: userInfoRequest.data.address,
        username: userInfoRequest.data.username,
      });
    }
  }, [userInfoRequest.loading, userInfoRequest.data]);

  React.useEffect(() => {
    if (!loginRequest.loading) {
      if (loginRequest.data?.status === "success") {
        userInfoRequest.fire();
        setError(undefined);
      }
      else if (!loginRequest.data || !loginRequest.ok) {
        // setError("登录失败");
        // alert("密码或者账户名错误");
      }
      setType("not defined");
    }
  }, [loginRequest.loading, loginRequest.data, loginRequest.ok]);

  const login = () => {
    loginRequest.fire({
      username: name,
      password: password
    });
    setType("login");
  };

  const Content = <Stack>
    <Stack.Item>
      <Text variant="xxLarge">登录账户</Text>
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="用户名" defaultValue={name} onChange={(_, v) => setName(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <TextField label="密码" canRevealPassword={true} type="password" defaultValue={password} onChange={(_, v) => setPassword(v)} />
    </Stack.Item>
    <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
      <Stack horizontal>
        <Stack.Item>
          <PrimaryButton text="登录" onClick={login} />
        </Stack.Item>
        <Stack.Item styles={{ root: { paddingLeft: 10 } }}>
          <DefaultButton text="注册" onClick={() => setType("register")} />
        </Stack.Item>
      </Stack>
    </Stack.Item>
    {
      error && <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <p style={{ color: "red" }}>{error}</p>
      </Stack.Item>
    }
  </Stack >;

  return type === "register" ? <Register setType={setType} /> : Content;
};

export default Login;