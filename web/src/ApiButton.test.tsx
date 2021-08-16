import { Icon } from "@material-ui/core"
import { mount } from "enzyme"
import React from "react"
import { ApiButton, ButtonText } from "./ApiButton"

type UIButton = Proto.v1alpha1UIButton

const ButtonName = "MyButton"

function newButton(icon?: string, text?: string): UIButton {
  return {
    metadata: {
      name: "button",
    },
    spec: {
      iconName: icon,
      text: text,
    },
  }
}

type Case = [
  showText: boolean,
  text?: string,
  iconName?: string,
  expectedIcon?: string,
  expectedText?: string
]

const cases: Case[] = [
  [true, "yum!", "cake", "yum!", "cake"],
  [true, "yum!", undefined, "yum!", undefined],
  // no text specified, fill in some default text ("Button")
  [true, undefined, "cake", "Button", "cake"],
  // ditto
  [true, undefined, undefined, "Button", undefined],
  [false, "yum!", "cake", undefined, "cake"],
  // showText=false and no icon specified, so use a default icon
  [false, "yum!", undefined, undefined, "smart_button"],
  [false, undefined, "cake", undefined, "cake"],
  // showText=false and no icon specified, so use a default icon
  [false, undefined, undefined, undefined, "smart_button"],
]

test.each(cases)(
  "button with showText %p, text %p, iconName %p",
  (showText, text, iconName, expectedText, expectedIcon) => {
    const button = newButton(iconName, text)
    const root = mount(<ApiButton button={button} showText={showText} />)

    const actualIcon = root.find(Icon)
    const expectedIconArray = expectedIcon ? [expectedIcon] : []
    expect(actualIcon.map((e) => e.text())).toEqual(expectedIconArray)

    const actualText = root.find(ButtonText)
    const expectedTextArray = expectedText ? [expectedText] : []
    expect(actualText.map((e) => e.text())).toEqual(expectedTextArray)
  }
)
