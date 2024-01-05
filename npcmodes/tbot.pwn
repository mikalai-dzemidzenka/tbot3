#include <a_npc>

new rec[64];
new recType;
new isSingle;
new groupID;

public OnRecordingPlaybackEnd() {
	if (isSingle == 0) {
		StartRecordingPlayback(recType,rec);
	}
}

public OnNPCSpawn(){
	SendCommand("/tbinit");
}

public OnClientMessage(color, text[]){
	switch (color) {
		case 0x000001: {
		    new idx;
			rec = strtok(text,idx);
			recType = strval(strtok(text,idx));
			isSingle = strval(strtok(text,idx));
			groupID = strval(strtok(text,idx));
			StartRecordingPlayback(recType,rec);
			SendCommand("/tbready");
			PauseRecordingPlayback();
			return;
		}
		
		case 0x000002: {
		    //sleep(1000);
			ResumeRecordingPlayback();
			return;
		}
	}
}


strtok(const string[], &index)
{
	new length = strlen(string);
	while ((index < length) && (string[index] <= ' '))
	{
		index++;
	}

	new offset = index;
	new result[32];
	while ((index < length) && (string[index] > ' ') && ((index - offset) < (sizeof(result) - 1)))
	{
		result[index - offset] = string[index];
		index++;
	}
	result[index - offset] = EOS;
	return result;
}
stock split(const strsrc[], strdest[][], delimiter)
{
    new i, li;
    new aNum;
    new len;
    while(i <= strlen(strsrc))
    {
        if(strsrc[i] == delimiter || i == strlen(strsrc))
        {
            len = strmid(strdest[aNum], strsrc, li, i, 128);
            strdest[aNum][len] = 0;
            li = i+1;
            aNum++;
        }
        i++;
    }
    return 1;
}
