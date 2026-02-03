# Ace
## Desciption 
Is a CLI quiz game to prep for interviews. Unlike leetcode where it tests on a more practical level, Ace tests you on the ability to answer interview questions with clarity, simplicity and fast. Don't get put on the spot the next time an interviewer asks you what an abstract class is, or what is a TCP. You might even know the answer, but if you can articulate it properly it might lose you the spot! And this is where this little tool comes in handy. 

## ID Formating 
To keep ids unique and consistant as the number of packs grow, ids are formated in this way:
**Pack ID**:

pack-[packhash]

*e.g,* `pack-packabc`

**Question ID**:

q-[packhash]-[TypeChoice]-[index]-[questionhash]

*e.g,* `q-a3f91c-choice-04-b91e2a`

## Rules 
- Do not directly change any files, use the CLI or TUI to do so 

