import * as React from "react";
import AppContext from "../../AppContext";
import { PrimaryButton, Stack, Text, TextField,Dropdown,IDropdownStyles, IDropdownOption,Label } from "@fluentui/react";

const postUrl = "/api/user/trans";

const User1: React.FunctionComponent = () => {
  const { user, setSelectedKey } = React.useContext(AppContext);
  const [title, setTitle] = React.useState<string>();
  const [title2, setTitle2] = React.useState<string>();
  const [address, setAddress] = React.useState<string>();
  const [man, setMan] = React.useState<string>();
  const [getmoney,setGetmoney] = React.useState<string>();
  const [money, setMoney] = React.useState<string>();
  const [id, setid] = React.useState<string>()
  const [currTags, setCurrTags] = React.useState<string>();
  const [tags, setTags] = React.useState<string[]>([]);
  const [results, setresults] = React.useState<string[]>([]);
  const [results2, setresults2] = React.useState<string[]>([]);
  const [text, setText] = React.useState<string>();
  const [type, setType] = React.useState<string>();
  const [selectedItem, setSelectedItem] = React.useState<IDropdownOption>();

  React.useEffect(() => {
    setSelectedKey && setSelectedKey("post");
  }, [])

  React.useEffect(() => {
    // console.log("title:",title);
    console.log("id:", id);
    // console.log("text:",text);
  }, [tags]);

  const post = () => {
    fetch(postUrl, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ funcName:"transfer_bill", funcParam: [title,title2,address,money], text: text })
    }).then(res => res.json()).then(data => {
      //console.log(data);
      //console.log(data.result);
      console.log(data.result[0]);
    //   setresults(data.result[]);
    //   setTags([]);
    //   setType(data.result);
    });
  };
  const post2 = () => {
    fetch(postUrl, {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ funcName: "settle", funcParam: [man,getmoney], text: text })
    }).then(res => res.json()).then(data => {
      //console.log(data);
      //console.log(data.result);
      console.log(data.result);
    //   setMoney(data.result[0].data);
    //   setTags([]);
    //   setType(data.result);
    });
  };

  const dropdownStyles: Partial<IDropdownStyles> = {
    dropdown: { width: 300 },
  };

  const onChange = (event: React.FormEvent<HTMLDivElement>, item: IDropdownOption): void => {
    setSelectedItem(item);
    setTitle(String(item.key));
  };

  const onChange2 = (event: React.FormEvent<HTMLDivElement>, item: IDropdownOption): void => {
    setSelectedItem(item);
    setTitle2(String(item.key));
  };


  const BeforePost =
    <Stack>
      <Stack.Item>
        <Text variant="xxLarge">个人金融操作</Text>
      </Stack.Item>
      <br></br>
      <Stack.Item>
        <Text variant="large">应收账款转移</Text>
      </Stack.Item>
      <Label>将自己持有的收据里的钱转移到自己的欠条里。举例：你欠A100元(记录该关系的欠条id为1)，B欠你200元(记录该关系的收据id为2)。transfer_bill(2,1,A的地址，100)，就会让B欠你100，你不欠A的钱</Label>
      {/* <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <TextField label="欠条id" onChange={(_, v) => setTags([...tags, v!])} />
      </Stack.Item> */}
      {/* <Dropdown
        placeholder="选择查看类型"
        label="选择需要查看的类型"
        selectedKey={selectedItem ? selectedItem.key : undefined}
        // defaultSelectedKeys={['apple', 'banana', 'grape']}
        options={options}
        styles={dropdownStyles}
        onChange={onChange}
      /> */}
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <TextField label="收据id" onChange={(_, v) => setTitle(v)} />
      </Stack.Item>
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <TextField label="欠条id" onChange={(_, v) => setTitle2(v)} />
      </Stack.Item>
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <TextField label="债权人地址" onChange={(_, v) => setAddress(v)} />
      </Stack.Item>
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <TextField label="转移金额" onChange={(_, v) => setMoney(v)} />
      </Stack.Item>
      {/* <Stack.Item styles={{ root: { paddingTop: 10, width: 1000 } }}>
        <Stack horizontal>
          <Stack.Item>
            <TextField value={currTags} label="ID" onChange={(_, v) => setCurrTags(v)} />
          </Stack.Item>
          <Stack.Item styles={{ root: { paddingLeft: 10, paddingTop: 30 } }}>
            <PrimaryButton text="确认ID" onClick={() => {
              if (currTags) {
                setTags([...tags, currTags!]);
                setCurrTags("");
              }
            }} />
          </Stack.Item>
          <Stack.Item styles={{ root: { paddingTop: 40, paddingLeft: 20 } }}>
            <i>当前ID：</i>
            {
              tags.map(
                (tag, index) => {
                  return <i key={index}>{" " + tag}</i>;
                }
              )
            }
          </Stack.Item>
        </Stack>
      </Stack.Item> */}
      
      {/* <Stack.Item styles={{ root: { paddingTop: 20, paddingLeft: 0 } }}>
            <Label>返回结果：</Label>
            {
              results.map(
                (result, index) => {
                  return <i key={index}>{" " + JSON.stringify(result)}</i>;
                }
              )
            }
      </Stack.Item> */}
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <PrimaryButton text="债务转移" onClick={post} />
      </Stack.Item>
      {/* <Stack.Item styles={{ root: { paddingTop: 10, width: 600, height: 400 } }}>
        <TextField label="正文" multiline rows={20} onChange={(_, v) => setText(v)} />
      </Stack.Item> */}
      <br></br>
      <br></br>
      <Stack.Item>
        <Text variant="large">结算(还钱）</Text>
      </Stack.Item>
      <Label>从自身资产中扣除对应欠条的数额并删去欠条</Label>
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <TextField label="债权人地址" onChange={(_, v) => setMan(v)} />
      </Stack.Item>
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <TextField label="欠条id" onChange={(_, v) => setGetmoney(v)} />
      </Stack.Item>
      
      {/* <Dropdown
        placeholder="选择查看类型"
        label="选择需要查看的类型"
        selectedKey={selectedItem ? selectedItem.key : undefined}
        // defaultSelectedKeys={['apple', 'banana', 'grape']}
        options={options2}
        styles={dropdownStyles}
        onChange={onChange2}
      /> */}
      
      {/* <Label>1、bank_info，获取银行信息。不需要参数。返回值包含两项，含义分别为：(1)、中央银行的名字。(2)中央银行的地址 </Label>
      <Label>2、register_company，注册个人公司，需注意个人最多只能有一个公司，多次使用会覆盖之前的结果。需要两个参数，分别为公司的名字和初始资金。无返回值</Label>
      <Label>3、getAsset，获取公司总资产。不需要参数。返回值包括一项，就是公司的总资产</Label>
      <Label>4、getDebt，根据欠条 id 获取欠条。需要一个参数即欠条id，欠条id从0开始。返回值有两项，分别为债权人的address，以及欠他的数额</Label>
      <Label>5、getReceipt，根据收据 id 获取收据。需要一个参数即收据id，欠条id从0开始。返回值有两项，分别为债务人的address，以及他欠的数额</Label>
      <Label>6、get_total_debt，获取一个公司总的应还账款（总负债额）。不需要参数。返回值有一项，为总负债数额</Label>
      <Label>7、get_total_receipt，获取一个公司总的应收账款（总放债额）。不需要参数。返回值有一项，为总放债数额</Label>
      <Label>8、create_bill，签发账单，只有银行账户可以操作。需要三个参数，分别为(1)债权人的地址(2)债务人的地址(3)借款数额。无返回值</Label>
      <Label>9、transfer_bill，应收账款转移，将自己持有的收据里的钱转移到自己的欠条里。需要四个参数，分别为(1)收据的id(2)欠条的id(3)欠条的债权人(4)转移的数值。无返回值。</Label>
      <Label>transfer_bill使用举例，你欠A100元(记录该关系的欠条id为1)，B欠你200元(记录该关系的收据id为2)。transfer_bill(2,1,A的地址，100)，就会让B欠你100，你不欠A的钱</Label>
      <Label>10、finance，融资，只有银行账户可以操作。需要两个参数即融资者的地址和融资额度(需小于融资者的总放债额度)，融资后公司资产增加相应数值。无返回值</Label>
      <Label>11、settle，结算(还钱），从自身资产中扣除对应欠条的数额并删去欠条。需要两个参数即还款人地址和欠条id。无返回值。</Label> */}
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <PrimaryButton text="结算" onClick={post2} />
      </Stack.Item>
      {
        (type === "not login" || type === "failure") && <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
          <p style={{ color: "red" }}>{type}</p>
        </Stack.Item>
      }
    </Stack>;

  const AfterPost =
    <Stack>
      <Stack.Item styles={{ root: { paddingTop: 20 } }}>
        <Text variant="xLarge">发布成功</Text>
      </Stack.Item>
      <Stack.Item styles={{ root: { paddingTop: 10, width: 300 } }}>
        <PrimaryButton text="再写一篇" onClick={() => {
          setType("initial");
          setCurrTags("");
          setTags([]);
        }} />
      </Stack.Item>
    </Stack>;


  return type === "success" ? AfterPost : BeforePost;
};

export default User1;
