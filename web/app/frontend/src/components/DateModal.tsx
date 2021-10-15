import React, {useState} from 'react';
import {Button} from "react-native-paper";
import {Platform} from "react-native";
import {DatePickerModal} from "react-native-paper-dates";

export function DateModal(props: { date: Date; setDate: React.Dispatch<React.SetStateAction<Date>>; }) {
    const [date, setDate] = props.date != null && props.setDate != null ? [props.date, props.setDate] : useState(new Date(Date.now()));
    // const [date, setDate] = useState(new Date(Date.now()));
    const [isOpen, setIsOpen] = useState(false);

    const show = () => setIsOpen(true)
    const hide = () => setIsOpen(false)

    let datePicker;

    if (Platform.OS === "web") {
        datePicker = (
            <input type="date"
                   value={toDateString(date)}
                   onChange={event => setDate(toDate(event.target.value))}/>
        )
    } else {
        datePicker = (
            <>
                <Button onPress={show}>{toDateString(date)}</Button>
                <DatePickerModal mode={"single"} onConfirm={(p) => {
                    hide();
                    setDate(p.date as Date)
                }} onDismiss={hide} visible={isOpen}/>
            </>
        )
    }

    return (
        <>
            {datePicker}
        </>
    )
}

function toDateString(date: Date): string {
    const y = date.getFullYear();
    const zeroPad = (n: number) => n < 9 ? '0' + (n + 1) : n + 1;
    const m = zeroPad(date.getMonth())
    const d = date.getDate();
    return `${y}-${m}-${d}`
}

function toDate(yyyymmdd: string): Date {
    const [yyyy, mm, dd] = yyyymmdd.split("-").map(s => s[0] === "0" ? +s[1] : +s);
    return new Date(yyyy, mm - 1, dd);
}
